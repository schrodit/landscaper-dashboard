// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package config

// Configuration defines the configuration for the dashboard server
type Configuration struct {
	HTTPPort  int `json:"httpPort"`
	HTTPSPort int `json:"httpsPort"`
}

// Default defaults the configuration values.
func Default(cfg *Configuration) {
	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = 8080
	}
}
