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
	"github.com/samuelngs/smartdns/net/http"
)

type httpServer struct {
	conf     *config.Config
	listener net.Listener
	stopped  bool
}

func (h *httpServer) listen() error {
	if h.conf.SNIProxy.HTTP.Enabled {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", h.conf.SNIProxy.HTTP.Port))
		if err != nil {
			logger.Warn(
				"could not listen for HTTP connections",
				log.String("error", err.Error()))
			return err
		}
		h.stopped = false
		h.listener = l
		h.acceptConnection(l)
	}
	return nil
}

func (h *httpServer) shutdown() {
	if !h.stopped {
		logger.Debug("stopped accepting HTTP connections")
		h.stopped = true
		h.listener.Close()
	}
}

func (h *httpServer) acceptConnection(l net.Listener) {
	logger.Debug("accepting HTTP connections", log.String("addr", l.Addr().String()))
	for {
		c, err := l.Accept()
		if err != nil {
			logger.Warn(
				"could not accept HTTP connection",
				log.String("error", err.Error()))
			if !h.stopped {
				continue
			}
			return
		}
		go h.handleConnection(c.(*net.TCPConn))
	}
}

func (h *httpServer) handleConnection(c *net.TCPConn) {
	defer c.Close()
	defer logger.Trace(
		"connection closed",
		log.String("remote-addr", c.RemoteAddr().String()),
		log.String("protocol", "http"))

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
		log.String("protocol", "http"))

	c.SetDeadline(time.Now().Add(h.conf.SNIProxy.ConnTimeout))

	hostname, prefix, err := http.ParseHost(c)
	if err != nil {
		logger.Warn(err.Error(), log.String("remote-addr", c.RemoteAddr().String()))
		return
	}

	logger.Trace("proxying http connection",
		log.String("remote-addr", c.RemoteAddr().String()),
		log.String("hostname", hostname))

	uri := net.JoinHostPort(hostname, "http")
	dst, err := net.DialTimeout("tcp", uri, h.conf.SNIProxy.DialTimeout)
	if err != nil {
		logger.Warn(
			"could not forward request",
			log.String("error", err.Error()),
			log.String("remote-addr", c.RemoteAddr().String()))
		return
	}

	if err := proxy(c, dst.(*net.TCPConn), h.conf.SNIProxy.DataTimeout, prefix); err != nil {
		logger.Warn(
			"could not proxy http connection",
			log.String("error", err.Error()),
			log.String("remote-addr", c.RemoteAddr().String()),
			log.String("hostname", hostname))
		return
	}
}
