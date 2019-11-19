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
	"fmt"
	"net"
	"time"

	"github.com/samuelngs/smartdns/log"
	"github.com/samuelngs/smartdns/net/http"
)

func (p *SNIProxy) listenHTTP(port int) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal(
			"could not listen for HTTP connections",
			log.String("error", err.Error()))
	}
	defer l.Close()

	logger.Debug("accepting HTTP connections", log.Int("port", port))
	for {
		c, err := l.Accept()
		if err != nil {
			logger.Fatal(
				"could not accept HTTP connection",
				log.String("error", err.Error()),
				log.String("remote-addr", c.RemoteAddr().String()))
		}
		if c, ok := c.(*net.TCPConn); ok {
			go p.handleHTTPConnection(c)
		}
	}
}

func (p *SNIProxy) handleHTTPConnection(c *net.TCPConn) {
	defer c.Close()
	defer logger.Trace(
		"connection closed",
		log.String("remote-addr", c.RemoteAddr().String()),
		log.String("protocol", "http"))

	logger.Trace(
		"connection accepted",
		log.String("remote-addr", c.RemoteAddr().String()),
		log.String("protocol", "http"))

	c.SetDeadline(time.Now().Add(p.connTimeout))

	hostname, prefix, err := http.ParseHost(c)
	if err != nil {
		logger.Warn(err.Error(), log.String("remote-addr", c.RemoteAddr().String()))
		return
	}

	logger.Trace("proxying http connection",
		log.String("remote-addr", c.RemoteAddr().String()),
		log.String("hostname", hostname))

	uri := net.JoinHostPort(hostname, "http")
	dst, err := net.DialTimeout("tcp", uri, p.dialTimeout)
	if err != nil {
		logger.Warn(
			"could not forward request",
			log.String("error", err.Error()),
			log.String("remote-addr", c.RemoteAddr().String()))
		return
	}

	if err := p.proxy(c, dst.(*net.TCPConn), prefix); err != nil {
		logger.Warn(
			"could not proxy http connection",
			log.String("error", err.Error()),
			log.String("remote-addr", c.RemoteAddr().String()),
			log.String("hostname", hostname))
		return
	}
}
