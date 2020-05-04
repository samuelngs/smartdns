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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/samuelngs/smartdns/net/ip"
)

// SNIProxy configuration
type SNIProxy struct {
	Host        string        `yaml:"host"`
	Ports       []string      `yaml:"ports"`
	ConnTimeout time.Duration `yaml:"conn_timeout"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	DataTimeout time.Duration `yaml:"data_timeout"`
}

// AllowedPorts returns all the open ports for server
func (s *SNIProxy) AllowedPorts() []int {
	portsmap := map[int]struct{}{}
	for _, rule := range s.Ports {
		switch ranges := strings.Split(rule, "-"); {
		case len(ranges) == 2:
			start, err := strconv.Atoi(ranges[0])
			if err != nil {
				continue
			}

			end, err := strconv.Atoi(ranges[1])
			if err != nil {
				continue
			}

			if start < 0 || end < 0 || end < start {
				continue
			}

			for i := start; i <= end; i++ {
				portsmap[i] = struct{}{}
			}
		case len(ranges) == 1:
			port, err := strconv.Atoi(ranges[0])
			if err != nil {
				continue
			}
			portsmap[port] = struct{}{}
		}
	}
	ports := make([]int, 0)
	for port := range portsmap {
		switch port {
		case 22, 53:
		default:
			ports = append(ports, port)
		}
	}
	sort.Sort(sort.IntSlice(ports))
	return ports
}

// DefaultSNIProxy configuration
func DefaultSNIProxy() *SNIProxy {
	p := &SNIProxy{
		Host:        "127.0.0.1",
		Ports:       []string{"0-10000"},
		ConnTimeout: time.Second * 20,
		DialTimeout: time.Second * 10,
		DataTimeout: time.Second * 240,
	}
	if host, ok := ip.FromEnv(); ok {
		p.Host = host.String()
	} else if host, ok := ip.FromIface(ip.Public, ip.Private, ip.Loopback); ok {
		p.Host = host.String()
	}
	return p
}
