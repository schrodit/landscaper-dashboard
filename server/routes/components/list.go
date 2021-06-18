// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package components

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/schrodit/landscaper-dashboard/server/routes"

	ocicd "github.com/gardener/component-spec/bindings-go/oci"
)

type ListComponentsRequest struct {
	RepositoryContext string `json:"repositoryContext"`
}

type ListComponentsResponse struct {
	RepositoryContext string   `json:"repositoryContext"`
	Components        []string `json:"components"`
}

type ListComponentVersionsRequest struct {
	RepositoryContext string `json:"repositoryContext"`
	ComponentName     string `json:"componentName"`
}

type ListComponentVersionsResponse struct {
	RepositoryContext string   `json:"repositoryContext"`
	ComponentName     string   `json:"componentName"`
	Versions          []string `json:"versions"`
}

type Component struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Ref     string `json:"ref"`
}

func (r *Router) ListComponents(data []byte) (*routes.Response, error) {
	req := ListComponentsRequest{}
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	if len(req.RepositoryContext) == 0 {
		return &routes.Response{
			Code: http.StatusBadRequest,
			Data: "a repository context has to be set",
		}, nil
	}
	ctx := context.Background()
	defer ctx.Done()

	componentsPrefix := path.Join(req.RepositoryContext, ocicd.ComponentDescriptorNamespace)
	repositories, err := r.ociClient.ListRepositories(ctx, componentsPrefix)
	if err != nil {
		return nil, err
	}

	res := ListComponentsResponse{}
	res.RepositoryContext = req.RepositoryContext

	res.Components = []string{}
	for _, repo := range repositories {
		res.Components = append(res.Components, strings.TrimPrefix(repo + "/", componentsPrefix))
	}
	return routes.OkResponse(res), nil
}

func (r *Router) ListComponentVersions(data []byte) (*routes.Response, error) {
	req := ListComponentVersionsResponse{}
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	if len(req.RepositoryContext) == 0 {
		return &routes.Response{
			Code: http.StatusBadRequest,
			Data: "a repository context has to be set",
		}, nil
	}
	if len(req.ComponentName) == 0 {
		return &routes.Response{
			Code: http.StatusBadRequest,
			Data: "a component name has to be set",
		}, nil
	}
	ctx := context.Background()
	defer ctx.Done()

	ref := path.Join(req.RepositoryContext, ocicd.ComponentDescriptorNamespace, req.ComponentName)
	tags, err := r.ociClient.ListTags(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("unable to list tags for %q: %w", ref, err)
	}

	res := ListComponentVersionsResponse{}
	res.RepositoryContext = req.RepositoryContext

	res.Versions = []string{}
	for _, tag := range tags {
		res.Versions = append(res.Versions, tag)
	}
	return routes.OkResponse(res), nil
}
