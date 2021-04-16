// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/schrodit/landscaper-dashboard/cmd/landscaper-dashboard/app"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	ctx := signals.SetupSignalHandler()
	if err := app.NewDashboardCmd(ctx).Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
