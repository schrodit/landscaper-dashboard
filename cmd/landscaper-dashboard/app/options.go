// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"flag"
	"io/ioutil"

	"github.com/go-logr/logr"
	"github.com/spf13/pflag"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/yaml"

	ociopts "github.com/gardener/component-cli/ociclient/options"

	"github.com/schrodit/landscaper-dashboard/server/config"
	"github.com/schrodit/landscaper-dashboard/server/logger"
)

type Options struct {
	ConfigPath  string
	FrontendDir string

	Log    logr.Logger
	Config config.Configuration
	// OciOptions contains all exposed options to configure the oci client.
	OciOptions ociopts.Options
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		fs = pflag.CommandLine
	}
	fs.StringVar(&o.ConfigPath, "config", "", "path to the configuration file")
	fs.StringVar(&o.FrontendDir, "frontend", "frontend/build", "path to the frontend directory")
	o.OciOptions.AddFlags(fs)
	logger.InitFlags(fs)
	fs.AddGoFlagSet(flag.CommandLine)
}

func (o *Options) Complete() error {
	log, err := logger.New(nil)
	if err != nil {
		return err
	}
	o.Log = log.WithName("setup")
	logger.SetLogger(log)
	ctrl.SetLogger(log)

	if len(o.ConfigPath) != 0 {
		data, err := ioutil.ReadFile(o.ConfigPath)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(data, &o.Config); err != nil {
			return err
		}
	}

	config.Default(&o.Config)
	if len(o.FrontendDir) != 0 {
		o.Config.FrontendDir = o.FrontendDir
	}
	return nil
}
