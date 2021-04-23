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
	Name      string
	Namespace string
	UpToDate  bool
	Phase     lsv1alpha1.ComponentInstallationPhase
}

func NewInstallationData(inst lsv1alpha1.Installation) InstallationData {
	return InstallationData{
		Name:      inst.Name,
		Namespace: inst.Namespace,
		UpToDate:  inst.Status.ObservedGeneration == inst.Generation,
		Phase:     inst.Status.Phase,
	}
}

func instListToInstDataList(insts lsv1alpha1.InstallationList) []InstallationData {
	res := make([]InstallationData, len(insts.Items))
	for i, elem := range insts.Items {
		res[i] = NewInstallationData(elem)
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
	instData := instListToInstDataList(*insts)

	return routes.OkResponse(instData), nil
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
