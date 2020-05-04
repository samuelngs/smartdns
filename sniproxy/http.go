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
	"io"
	"net"
	"time"

	"github.com/samuelngs/smartdns/config"
	"github.com/samuelngs/smartdns/log"
	"github.com/samuelngs/smartdns/net/http"
	"github.com/samuelngs/smartdns/net/https"
)

type httpServer struct {
	conf     *config.Config
	port     int
	listener net.Listener
	started  bool
}

func (h *httpServer) listen() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", h.port))
	if err != nil {
		logger.Warn(
			"could not listen for HTTP and HTTPS connections",
			log.String("error", err.Error()))
		return err
	}
	h.started = false
	h.listener = l
	h.acceptConnection(l)
	return nil
}

func (h *httpServer) shutdown() {
	if h.started {
		logger.Debug("stopped accepting HTTP and HTTPS connections")
		h.started = true
		h.listener.Close()
	}
}

func (h *httpServer) acceptConnection(l net.Listener) {
	logger.Debug("accepting HTTP and HTTPS connections", log.String("addr", l.Addr().String()))
	for {
		c, err := l.Accept()
		if err != nil {
			logger.Warn(
				"could not accept HTTP or HTTPS connection",
				log.String("error", err.Error()))
			if h.started {
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
		log.String("remote-addr", c.RemoteAddr().String()))

	if !h.conf.Network.IsAllowedIP(c.RemoteAddr()) {
		logger.Trace(
			"connection rejected",
			log.String("remote-addr", c.RemoteAddr().String()))
		return
	}

	logger.Trace(
		"connection accepted",
		log.String("remote-addr", c.RemoteAddr().String()))

	c.SetDeadline(time.Now().Add(h.conf.SNIProxy.ConnTimeout))

	logger.Trace(
		"checking connection protocol",
		log.String("remote-addr", c.RemoteAddr().String()))

	f := make([]byte, 1)
	c.Read(f)

	if f[0] == 22 {
		h.handleHTTPSConnection(c)
		return
	}

	hostname, prefix, err := http.ParseHost(c, f)
	if err != nil {
		logger.Warn(err.Error(), log.String("remote-addr", c.RemoteAddr().String()))
		return
	}

	h.handleHTTPConnection(c, hostname, prefix)
}

func (h *httpServer) handleHTTPConnection(c *net.TCPConn, hostname string, prefix io.Reader) {
	logger.Trace("proxying http connection",
		log.String("remote-addr", c.RemoteAddr().String()),
		log.String("hostname", hostname))

	uri := net.JoinHostPort(hostname, fmt.Sprintf("%d", h.port))
	dst, err := net.DialTimeout("tcp", uri, h.conf.SNIProxy.DialTimeout)
	if err != nil {
		logger.Warn(
			"could not forward http request",
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

func (h *httpServer) handleHTTPSConnection(c *net.TCPConn) {
	logger.Trace("reading sni-hostname",
		log.String("remote-addr", c.RemoteAddr().String()))

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

	logger.Trace("proxying https connection",
		log.String("remote-addr", c.RemoteAddr().String()),
		log.String("hostname", m.Hostname))

	uri := net.JoinHostPort(m.Hostname, fmt.Sprintf("%d", h.port))
	dst, err := net.DialTimeout("tcp", uri, h.conf.SNIProxy.DialTimeout)
	if err != nil {
		logger.Warn(
			"could not forward https request",
			log.String("error", err.Error()),
			log.String("remote-addr", c.RemoteAddr().String()))
		return
	}

	if err := proxy(c, dst.(*net.TCPConn), h.conf.SNIProxy.DataTimeout, &m.Buffer); err != nil {
		logger.Warn(
			"could not proxy https connection",
			log.String("error", err.Error()),
			log.String("remote-addr", c.RemoteAddr().String()),
			log.String("hostname", m.Hostname))
		return
	}
}
