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
)

// Network configuration
type Network struct {
	AllowedIPs []net.IP `yaml:"allowed_ips"`
	BlockedIPs []net.IP `yaml:"blocked_ips"`
}

// DefaultNetwork configuration
func DefaultNetwork() *Network {
	return &Network{
		AllowedIPs: make([]net.IP, 0),
		BlockedIPs: make([]net.IP, 0),
	}
}

// IsAllowedIP checks if the ip is allowed to make requests to this server
func (n *Network) IsAllowedIP(s interface{}) bool {
	if len(n.AllowedIPs) == 0 && len(n.BlockedIPs) == 0 {
		return true
	}
	var ip net.IP
	switch o := s.(type) {
	case string:
		ip = net.ParseIP(o)
	case net.IP:
		ip = o
	case *net.IPNet:
		ip = o.IP
	case *net.IPAddr:
		ip = o.IP
	case *net.UDPAddr:
		ip = o.IP
	case *net.TCPAddr:
		ip = o.IP
	}
	if ip == nil {
		return false
	}

	if len(n.BlockedIPs) > 0 {
		for _, bip := range n.BlockedIPs {
			if bip.Equal(ip) {
				return false
			}
		}
	}
	if len(n.AllowedIPs) > 0 {
		for _, aip := range n.AllowedIPs {
			if aip.Equal(ip) {
				return true
			}
		}
		return false
	}

	return true
}
