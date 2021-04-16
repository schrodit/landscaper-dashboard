// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"flag"
	"io/ioutil"

	"github.com/go-logr/logr"
	"github.com/schrodit/landscaper-dashboard/server/config"
	"github.com/schrodit/landscaper-dashboard/server/logger"
	"github.com/spf13/pflag"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/yaml"
)

type Options struct {
	ConfigPath string

	Log    logr.Logger
	Config config.Configuration
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		fs = pflag.CommandLine
	}
	fs.StringVar(&o.ConfigPath, "config", "", "path to the configuration file")
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
	return nil
}
