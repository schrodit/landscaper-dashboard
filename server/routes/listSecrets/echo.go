// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package lsRouter

import (
	"context"

	"github.com/schrodit/landscaper-dashboard/server/routes"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LSRouter struct {
	client client.Client
}

func NewLsRouter(c client.Client) *LSRouter {
	return &LSRouter{
		client: c,
	}
}

func (r LSRouter) ListSecrets(_ []byte) (*routes.Response, error) {
	ctx := context.Background()
	defer ctx.Done()
	secrets := &corev1.SecretList{}
	if err := r.client.List(ctx, secrets); err != nil {
		return nil, err
	}

	return routes.OkResponse(secrets), nil
}
