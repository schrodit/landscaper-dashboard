// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package components

import (
	"github.com/gardener/component-cli/ociclient"
	"github.com/go-logr/logr"

	"github.com/schrodit/landscaper-dashboard/server/routes"
)

// Router is the router that contains all component related api routes
type Router struct {
	log       logr.Logger
	ociClient ociclient.ExtendedClient
}

func New(log logr.Logger, ociClient ociclient.ExtendedClient) *Router {
	return &Router{
		log:       log,
		ociClient: ociClient,
	}
}

func AddToRouter(log logr.Logger, ociClient ociclient.ExtendedClient, router *routes.Router) {
	r := New(log, ociClient)
	router.AddCommonPath("listComponents", r.ListComponents)
	router.AddCommonPath("listComponentVersions", r.ListComponentVersions)
}
