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

package sniproxy

import (
	"github.com/samuelngs/smartdns/config"
)

// SNIProxy constructs a sni-proxy server
type SNIProxy struct {
	conf *config.Config
}

// Start initializes and starts sni-proxy server
func (p *SNIProxy) Start() error {
	if p.conf.Proxy.HTTP.Enabled {
		go p.listenHTTP(p.conf.Proxy.HTTP.Port)
	}
	if p.conf.Proxy.HTTPS.Enabled {
		go p.listenHTTPS(p.conf.Proxy.HTTPS.Port)
	}
	return nil
}

// Stop stops the running sni-proxy server
func (p *SNIProxy) Stop() error {
	return nil
}
