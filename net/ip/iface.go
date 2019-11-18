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

package ip

import (
	"net"
)

var localhost = net.ParseIP("127.0.0.1")

// AddressType presents the type of an network address
type AddressType int8

// defines a list of available network address types
const (
	None AddressType = iota
	Loopback
	Private
	Public
)

func fromIface(types ...AddressType) (net.IP, bool) {
	if len(types) == 0 {
		return nil, false
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, false
	}
	var ips []net.IP
	for _, i := range ifaces {
		if addrs, err := i.Addrs(); err == nil {
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				if ip != nil && !ip.Equal(localhost) {
					ips = append(ips, ip)
				}
			}
		}
	}
	for _, typ := range types {
		for _, ip := range ips {
			switch {
			case ip != nil && isPublicIP(ip) && typ == Public:
				return ip, true
			case ip != nil && isPrivateIP(ip) && typ == Private:
				return ip, true
			case ip != nil && ip.IsLoopback() && typ == Loopback:
				return ip, true
			}
		}
	}
	return localhost, false
}
