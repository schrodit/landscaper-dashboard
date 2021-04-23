// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package instData

import (
	"context"
	"encoding/json"
	"fmt"

	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/schrodit/landscaper-dashboard/server/routes"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type InstallationRouter struct {
	client client.Client
	logger logr.Logger
}

type options struct {
	Namespace               string
	IncludeSubinstallations bool
}

func newOptions() *options {
	return &options{
		Namespace:               "",
		IncludeSubinstallations: false,
	}
}

func (o *options) parseFromBytes(data []byte) error {
	err := json.Unmarshal(data, o)
	return err
}

type InstallationData struct {
	Name             string
	Namespace        string
	UpToDate         bool
	Phase            lsv1alpha1.ComponentInstallationPhase
	Subinstallations []InstallationData
	Execution        *ExecutionData
}

type ExecutionData struct {
	Name        string
	Namespace   string
	UpToDate    bool
	Phase       lsv1alpha1.ExecutionPhase
	DeployItems []DeployItemData
}

type DeployItemData struct {
	Name      string
	Namespace string
	UpToDate  bool
	Phase     lsv1alpha1.ExecutionPhase
}

func NewInstallationData(inst lsv1alpha1.Installation) *InstallationData {
	return &InstallationData{
		Name:      inst.Name,
		Namespace: inst.Namespace,
		UpToDate:  inst.Status.ObservedGeneration == inst.Generation,
		Phase:     inst.Status.Phase,
	}
}

func (id *InstallationData) SetExecutionData(exec lsv1alpha1.Execution) *ExecutionData {
	id.Execution = &ExecutionData{
		Name:      exec.Name,
		Namespace: exec.Namespace,
		UpToDate:  exec.Status.ObservedGeneration == exec.Generation,
		Phase:     exec.Status.Phase,
	}
	return id.Execution
}

func (ed *ExecutionData) SetDeployItemData(dis []lsv1alpha1.DeployItem) {
	ed.DeployItems = make([]DeployItemData, len(dis))
	for i, elem := range dis {
		ed.DeployItems[i] = DeployItemData{
			Name:      elem.Name,
			Namespace: elem.Namespace,
			UpToDate:  elem.Status.ObservedGeneration == elem.Generation,
			Phase:     elem.Status.Phase,
		}
	}
}

func (r *InstallationRouter) instListToInstDataList(insts lsv1alpha1.InstallationList, execs lsv1alpha1.ExecutionList, dis lsv1alpha1.DeployItemList) []*InstallationData {
	res := make([]*InstallationData, len(insts.Items))
	for i, elem := range insts.Items {
		var ex *lsv1alpha1.Execution
		if elem.Status.ExecutionReference != nil {
			for _, exec := range execs.Items {
				if exec.Name == elem.Status.ExecutionReference.Name && exec.Namespace == elem.Status.ExecutionReference.Namespace {
					ex = &exec
					break
				}
			}
		}
		res[i] = NewInstallationData(elem)
		if ex != nil {
			res[i].SetExecutionData(*ex)

			var exDis []lsv1alpha1.DeployItem
			if ex.Status.DeployItemReferences != nil {
				exDis = []lsv1alpha1.DeployItem{}
				for _, dr := range ex.Status.DeployItemReferences {
					for _, d := range dis.Items {
						if dr.Reference.Name == d.Name && dr.Reference.Namespace == d.Namespace {
							exDis = append(exDis, d)
							break
						}
					}
				}
				res[i].Execution.SetDeployItemData(exDis)
			}
		}
	}
	return res
}

func NewInstallationRouter(c client.Client, log logr.Logger) *InstallationRouter {
	return &InstallationRouter{
		client: c,
		logger: log,
	}
}

func (r *InstallationRouter) ListInstallations(data []byte) (*routes.Response, error) {
	ctx := context.Background()
	defer ctx.Done()
	o := newOptions()
	if len(data) > 0 {
		r.logger.Info("parsing received options", "options", string(data)) // todo V(7)
		err := o.parseFromBytes(data)
		if err != nil {
			r.logger.Error(err, "unable to unmarshal options", "data", string(data))
			return &routes.Response{Code: 500}, err
		}
	}

	insts, err := r.fetchInstallations(ctx, o)
	if err != nil {
		r.logger.Error(err, "error trying to fetch the installations")
		return &routes.Response{Code: 500}, err
	}
	execs, err := r.fetchExecutions(ctx, insts)
	if err != nil {
		r.logger.Error(err, "error trying to fetch the executions")
		return &routes.Response{Code: 500}, err
	}
	dis, err := r.fetchDeployItems(ctx, execs)
	if err != nil {
		r.logger.Error(err, "error trying to fetch the deploy items")
		return &routes.Response{Code: 500}, err
	}

	instData := r.instListToInstDataList(*insts, *execs, *dis)

	return routes.OkResponse(instData), nil
}

func (r *InstallationRouter) fetchInstallations(ctx context.Context, o *options) (*lsv1alpha1.InstallationList, error) {
	lopts, err := r.buildListOptions(ctx, o)
	if err != nil {
		return nil, fmt.Errorf("unable to build list options: %w", err)
	}
	insts := &lsv1alpha1.InstallationList{}
	if err := r.client.List(ctx, insts, lopts); err != nil {
		return nil, fmt.Errorf("unable to list installations: %w", err)
	}
	return insts, nil
}

func (r *InstallationRouter) fetchExecutions(ctx context.Context, insts *lsv1alpha1.InstallationList) (*lsv1alpha1.ExecutionList, error) {
	execs := &lsv1alpha1.ExecutionList{}
	if err := r.client.List(ctx, execs); err != nil {
		return nil, fmt.Errorf("unable to list executions: %w", err)
	}
	execNames := make(map[lsv1alpha1.ObjectReference]bool, len(insts.Items))
	for _, elem := range insts.Items {
		if elem.Status.ExecutionReference != nil {
			execNames[*elem.Status.ExecutionReference] = true
		}
	}
	return filterExecutionList(execs, func(e *lsv1alpha1.Execution) bool {
		_, ok := execNames[lsv1alpha1.ObjectReference{Name: e.Name, Namespace: e.Namespace}]
		return ok
	}), nil
}

func (r *InstallationRouter) fetchDeployItems(ctx context.Context, execs *lsv1alpha1.ExecutionList) (*lsv1alpha1.DeployItemList, error) {
	dis := &lsv1alpha1.DeployItemList{}
	if err := r.client.List(ctx, dis); err != nil {
		return nil, fmt.Errorf("unable to list deploy items: %w", err)
	}
	flatDIRefs := map[lsv1alpha1.ObjectReference]bool{}
	for _, ex := range execs.Items {
		if ex.Status.DeployItemReferences != nil {
			for _, diRef := range ex.Status.DeployItemReferences {
				flatDIRefs[lsv1alpha1.ObjectReference{Name: diRef.Reference.Name, Namespace: diRef.Reference.Namespace}] = true
			}
		}
	}
	return filterDeployItemList(dis, func(e *lsv1alpha1.DeployItem) bool {
		_, ok := flatDIRefs[lsv1alpha1.ObjectReference{Name: e.Name, Namespace: e.Namespace}]
		return ok
	}), nil
}

func (r *InstallationRouter) buildListOptions(ctx context.Context, o *options) (*client.ListOptions, error) {
	lopts := &client.ListOptions{}
	logger := r.logger.WithName("buildListOptions") // todo V(7)
	if len(o.Namespace) > 0 {
		logger.Info("namespace configured", "value", o.Namespace)
		lopts.Namespace = o.Namespace
	}
	logger.Info("include subinstallations", "value", o.IncludeSubinstallations)
	if !o.IncludeSubinstallations {
		req, err := labels.NewRequirement(lsv1alpha1.EncompassedByLabel, selection.DoesNotExist, []string{})
		if err != nil {
			return nil, fmt.Errorf("unable to build requirement: %w", err)
		}
		lopts.LabelSelector = labels.NewSelector().Add(*req)
	}
	return lopts, nil
}

func filterInstallationList(il *lsv1alpha1.InstallationList, accept func(*lsv1alpha1.Installation) bool) *lsv1alpha1.InstallationList {
	res := &lsv1alpha1.InstallationList{
		TypeMeta: il.TypeMeta,
		ListMeta: il.ListMeta,
		Items:    []lsv1alpha1.Installation{},
	}
	for _, elem := range il.Items {
		if accept(&elem) {
			res.Items = append(res.Items, elem)
		}
	}
	return res
}
func filterExecutionList(il *lsv1alpha1.ExecutionList, accept func(*lsv1alpha1.Execution) bool) *lsv1alpha1.ExecutionList {
	res := &lsv1alpha1.ExecutionList{
		TypeMeta: il.TypeMeta,
		ListMeta: il.ListMeta,
		Items:    []lsv1alpha1.Execution{},
	}
	for _, elem := range il.Items {
		if accept(&elem) {
			res.Items = append(res.Items, elem)
		}
	}
	return res
}
func filterDeployItemList(il *lsv1alpha1.DeployItemList, accept func(*lsv1alpha1.DeployItem) bool) *lsv1alpha1.DeployItemList {
	res := &lsv1alpha1.DeployItemList{
		TypeMeta: il.TypeMeta,
		ListMeta: il.ListMeta,
		Items:    []lsv1alpha1.DeployItem{},
	}
	for _, elem := range il.Items {
		if accept(&elem) {
			res.Items = append(res.Items, elem)
		}
	}
	return res
}
