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
	"fmt"
	"os"
)

// DNS configuration
type DNS struct {
	TLS            *DNSTLS       `yaml:"tls"`
	DNSResolveList []*DNSResolve `yaml:"resolve_dns"`
}

// DNSTLS configuration
type DNSTLS struct {
	Enabled  bool   `yaml:"enabled"`
	Email    string `yaml:"email"`
	Hostname string `yaml:"hostname"`
}

// DefaultDNS generates default settings for DNS
func DefaultDNS() *DNS {
	return &DNS{
		TLS:            DefaultDNSTLS(),
		DNSResolveList: make([]*DNSResolve, 0),
	}
}

// DefaultDNSTLS generates default settings for dns-tls
func DefaultDNSTLS() *DNSTLS {
	return &DNSTLS{
		Enabled:  true,
		Email:    fmt.Sprintf("admin@%s", os.Getenv("hostname")),
		Hostname: os.Getenv("hostname"),
	}
}
