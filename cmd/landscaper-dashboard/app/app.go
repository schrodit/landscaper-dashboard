// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gardener/landscaper/apis/core/install"
	"github.com/gin-gonic/gin"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/schrodit/landscaper-dashboard/server/routes"
	"github.com/schrodit/landscaper-dashboard/server/routes/components"
	"github.com/schrodit/landscaper-dashboard/server/routes/echo"
	"github.com/schrodit/landscaper-dashboard/server/routes/frontend"
	instData "github.com/schrodit/landscaper-dashboard/server/routes/listInstallationData"
	lsRouter "github.com/schrodit/landscaper-dashboard/server/routes/listSecrets"
	"github.com/spf13/cobra"
	"gopkg.in/olahol/melody.v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func NewDashboardCmd(ctx context.Context) *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use:          "landscaper-dashboard [--config]",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Complete(); err != nil {
				return err
			}
			return opts.Run(ctx)
		},
	}
	opts.AddFlags(cmd.Flags())
	return cmd
}

func (o *Options) Run(ctx context.Context) error {
	restConfig := ctrl.GetConfigOrDie()
	mgr, err := ctrl.NewManager(restConfig, manager.Options{
		MetricsBindAddress:     "0",
		HealthProbeBindAddress: "0",
	})
	if err != nil {
		return fmt.Errorf("unable to setup manager: %w", err)
	}
	install.Install(mgr.GetScheme())

	ociClient, _, err := o.OciOptions.Build(o.Log, osfs.New())
	if err != nil {
		return err
	}

	server := gin.Default()
	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", o.Config.HTTPPort),
		Handler: server,
	}
	m := melody.New() // websocket middleware

	if len(o.Config.FrontendDir) != 0 {
		if err := frontend.AddToRouter(o.Config.FrontendDir, server); err != nil {
			return err
		}
		o.Log.Info("successfully registered static frontend")
	}

	server.GET("/ws", func(c *gin.Context) {
		if err := m.HandleRequest(c.Writer, c.Request); err != nil {
			o.Log.Error(err, "unable to handle websocket request")
		}
	})

	router := routes.NewRouter(o.Log.WithName("router"), server)
	if err := router.Register(m); err != nil {
		return err
	}
	router.AddCommonPath("echo", echo.Route)
	router.AddCommonPath("listSecrets", lsRouter.NewLsRouter(mgr.GetClient()).ListSecrets)
	components.AddToRouter(o.Log, ociClient, router)
	router.AddCommonPath("listInstallationData", instData.NewInstallationRouter(mgr.GetClient(), o.Log.WithName("listInstallationData")).ListInstallations)
	go func() {
		if err := mgr.Start(ctx); err != nil {
			o.Log.Error(err, "unable to start controller manager")
			os.Exit(1)
		}
		if err := httpServer.Shutdown(ctx); err != nil {
			o.Log.Error(err, "unable to shutdown http server")
			os.Exit(1)
		}
	}()
	// todo: enable https
	o.Log.Info(fmt.Sprintf("Start http server listening on %d", o.Config.HTTPPort))
	return httpServer.ListenAndServe()
}
