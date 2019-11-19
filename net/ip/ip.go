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
	"os"
)

// FromEnv reads an network address from environment variable
func FromEnv() (net.IP, bool) {
	env, ok := os.LookupEnv("HOST")
	if !ok {
		return nil, false
	}
	ip := net.ParseIP(env)
	return ip, ip != nil
}

// FromIface reads an network address from system network interface
func FromIface(types ...AddressType) (net.IP, bool) {
	return fromIface(types...)
}

// FromEth0 reads public address of the server from eth0.me
func FromEth0() (net.IP, bool) {
	return fromExternalService("https://eth0.me")
}

// FromIPEcho reads public address of server from ipecho.net
func FromIPEcho() (net.IP, bool) {
	return fromExternalService("https://ipecho.net/plain")
}

// FromExternalService reads the public address of server from an external service, when the
// service receives a request from the sender, typically this server, it should return the
// client ip in the response body
//
// example:
// $ curl eth0.me => 43.54.65.76
// $ curl ipecho.net/plain => 43.54.65.76
//
func FromExternalService(addr string) (net.IP, bool) {
	return fromExternalService(addr)
}
