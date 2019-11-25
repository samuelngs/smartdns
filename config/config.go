// Copyright 2019 smartdns authors
// This file is part of the smartdns library.
//
// The smartdns library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The smartdns library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the smartdns library. If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

// Config encapsulates all configuration details for smartdns
type Config struct {
	Path     string    `yaml:"path"`
	Network  *Network  `yaml:"network"`
	DNS      *DNS      `yaml:"dns"`
	SNIProxy *SNIProxy `yaml:"proxy"`
}

// DefaultConfig generates the default settings for smartdns
func DefaultConfig() *Config {
	return &Config{
		Path:     os.ExpandEnv("/etc/smartdns/smartdns.yaml"),
		Network:  DefaultNetwork(),
		DNS:      DefaultDNS(),
		SNIProxy: DefaultSNIProxy(),
	}
}

// Read reads the yaml configuration from bytes
func Read(b []byte) (*Config, error) {
	config := DefaultConfig()
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return config, nil
}

// FromFile reads the yaml configuration from specify path
func FromFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf, err := Read(data)
	if err != nil {
		return nil, err
	}
	conf.Path = path
	return conf, nil
}
