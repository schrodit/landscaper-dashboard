// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package components

import (
	"encoding/json"

	"github.com/schrodit/landscaper-dashboard/server/routes"
)

type ListComponentsRequest struct {
	RepositoryContext string `json:"repositoryContext"`
}

type ListComponentsResponse struct {
	RepositoryContext string      `json:"repositoryContext"`
	Components        []Component `json:"components"`
}

type Component struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Ref     string `json:"ref"`
}

func ListComponents(data []byte) (*routes.Response, error) {
	req := ListComponentsRequest{}
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}
	res := ListComponentsResponse{}
	res.RepositoryContext = req.RepositoryContext
	res.Components = []Component{
		{
			Name:    "my-component",
			Version: "0.0.1",
			Ref:     "dockerhub.io/my-component:0.0.1",
		},
		{
			Name:    "my-other-component",
			Version: "0.0.2",
			Ref:     "dockerhub.io/my-other-component:0.0.2",
		},
	}
	return routes.OkResponse(res), nil
}
