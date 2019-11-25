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

	"github.com/samuelngs/smartdns/config"
	"github.com/samuelngs/smartdns/log"
	"github.com/samuelngs/smartdns/net/https"
)

type httpsServer struct {
	conf     *config.Config
	listener net.Listener
	stopped  bool
}

func (h *httpsServer) listen() error {
	if h.conf.SNIProxy.HTTPS.Enabled {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", h.conf.SNIProxy.HTTPS.Port))
		if err != nil {
			logger.Fatal(
				"could not listen for HTTPS connections",
				log.String("error", err.Error()))
			return err
		}
		h.stopped = false
		h.listener = l
		h.acceptConnection(l)
	}
	return nil
}

func (h *httpsServer) shutdown() {
	logger.Debug("stopped accepting HTTPS connections")
	if !h.stopped {
		h.stopped = true
		h.listener.Close()
	}
}

func (h *httpsServer) acceptConnection(l net.Listener) {
	logger.Debug("accepting HTTPS connections", log.String("addr", l.Addr().String()))
	for {
		c, err := l.Accept()
		if err != nil {
			logger.Warn(
				"could not accept HTTPS connection",
				log.String("error", err.Error()))
			if !h.stopped {
				continue
			}
			return
		}
		go h.handleConnection(c.(*net.TCPConn))
	}
}

func (h *httpsServer) handleConnection(c *net.TCPConn) {
	defer c.Close()
	defer logger.Trace(
		"connection closed",
		log.String("remote-addr", c.RemoteAddr().String()),
		log.String("protocol", "https"))

	if !h.conf.Network.IsAllowedIP(c.RemoteAddr()) {
		logger.Trace(
			"connection rejected",
			log.String("remote-addr", c.RemoteAddr().String()),
			log.String("protocol", "https"))
		return
	}

	logger.Trace(
		"connection accepted",
		log.String("remote-addr", c.RemoteAddr().String()),
		log.String("protocol", "https"))

	c.SetDeadline(time.Now().Add(h.conf.SNIProxy.ConnTimeout))

	m, err := https.ParseHandshakeMessage(c)
	if err != nil {
		logger.Warn(
			"could not read sni-hostname",
			log.String("remote-addr", c.RemoteAddr().String()),
			log.String("error", err.Error()))
		return
	}
	if len(m.Hostname) == 0 {
		logger.Warn(
			"could not read sni-hostname",
			log.String("remote-addr", c.RemoteAddr().String()))
		return
	}

	logger.Trace("proxying http connection",
		log.String("remote-addr", c.RemoteAddr().String()),
		log.String("hostname", m.Hostname))

	uri := net.JoinHostPort(m.Hostname, "https")
	dst, err := net.DialTimeout("tcp", uri, h.conf.SNIProxy.DialTimeout)
	if err != nil {
		logger.Warn(
			"could not forward request",
			log.String("error", err.Error()),
			log.String("remote-addr", c.RemoteAddr().String()))
		return
	}

	if err := proxy(c, dst.(*net.TCPConn), h.conf.SNIProxy.DataTimeout, &m.Buffer); err != nil {
		logger.Warn(
			"could not proxy http connection",
			log.String("error", err.Error()),
			log.String("remote-addr", c.RemoteAddr().String()),
			log.String("hostname", m.Hostname))
		return
	}
}
