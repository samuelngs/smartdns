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
	"net"
	"strings"
)

// DNSResolveList represents a list of dns resolver list
type DNSResolveList []*DNSResolve

// DNSResolve query rule
type DNSResolve struct {
	Name       string `yaml:"name"`
	Nameserver string `yaml:"nameserver,omitempty"`
	IP         string `yaml:"ip,omitempty"`
	TTL        int    `yaml:"ttl"`
}

// IsValid returns true if the custom dns configuration is valid
func (d *DNSResolve) IsValid() bool {
	switch {
	case len(d.Nameserver) > 0 && len(d.IP) > 0:
		return false
	case len(d.Nameserver) == 1 && d.Nameserver == "-":
		return true
	case len(d.Nameserver) > 0:
		return net.ParseIP(d.Nameserver) != nil
	case len(d.IP) > 0:
		return net.ParseIP(d.IP) != nil
	}
	return false
}

// NameserverAddr returns the address of nameserver
func (d *DNSResolve) NameserverAddr() string {
	if _, _, err := net.SplitHostPort(d.Nameserver); err != nil {
		return net.JoinHostPort(d.Nameserver, "53")
	}
	return d.Nameserver
}

// ResolveWithProxy resolves target with proxy server
func ResolveWithProxy(name string, ttl int) *DNSResolve {
	d := &DNSResolve{Name: name, Nameserver: "-", TTL: ttl}
	if d.IsValid() {
		return d
	}
	return nil
}

// ResolveToIP points a domain name to an ip address
func ResolveToIP(name, ip string, ttl int) *DNSResolve {
	d := &DNSResolve{Name: name, IP: ip, TTL: ttl}
	if d.IsValid() {
		return d
	}
	return nil
}

// ResolveWithNameserver sets the dns of domain name to be resolved
// through an external dns server
func ResolveWithNameserver(name, ns string, ttl int) *DNSResolve {
	d := &DNSResolve{Name: name, Nameserver: ns, TTL: ttl}
	if d.IsValid() {
		return d
	}
	return nil
}

// MatchDNS returns the dns rule that matches the hostname
func (d DNSResolveList) MatchDNS(name string) *DNSResolve {
	var matches []*DNSResolve
	for _, rule := range d {
		if rule != nil {
			o := rule.Name
			if len(o) > 0 && o[len(o)-1] != 0x2e {
				o += string(0x2e)
			}
			if strings.HasSuffix(name, o) {
				matches = append(matches, rule)
			}
		}
	}
	if len(matches) > 0 {
		match := matches[0]
		for _, rule := range matches {
			if len(rule.Name) > len(match.Name) {
				match = rule
			}
		}
		return match
	}
	return nil
}
