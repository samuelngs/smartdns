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
	"time"

	"github.com/samuelngs/smartdns/net/ip"
)

// Proxy configuration
type Proxy struct {
	Host        string        `yaml:"host"`
	ConnTimeout time.Duration `yaml:"conn_timeout"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	DataTimeout time.Duration `yaml:"data_timeout"`
	HTTP        *HTTPProxy    `yaml:"http"`
	HTTPS       *HTTPSProxy   `yaml:"https"`
	DNS         []*DNS        `yaml:"resolve_dns"`
}

// DefaultProxy configuration
func DefaultProxy() *Proxy {
	p := &Proxy{
		Host:        "127.0.0.1",
		ConnTimeout: time.Second * 20,
		DialTimeout: time.Second * 10,
		DataTimeout: time.Second * 240,
		HTTP:        DefaultHTTPProxy(),
		HTTPS:       DefaultHTTPSProxy(),
		DNS:         make([]*DNS, 0),
	}
	if host, ok := ip.FromEnv(); ok {
		p.Host = host.String()
	} else if host, ok := ip.FromIface(ip.Public, ip.Private, ip.Loopback); ok {
		p.Host = host.String()
	}
	return p
}
