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

package dnsproxy

import (
	"context"
	"crypto/tls"

	"github.com/miekg/dns"
	"github.com/samuelngs/smartdns/config"
	"github.com/samuelngs/smartdns/log"
	"golang.org/x/crypto/acme"
	"golang.org/x/sync/errgroup"
)

var logger = log.DefaultLogger

// DNSProxy constructs a dns-proxy server
type DNSProxy struct {
	conf   *config.Config
	acme   *acme.Client
	dns    *dnsServer
	dnstls *dnsServer
	ctx    context.Context
}

// Start initializes and starts dns-proxy server
func (d *DNSProxy) Start() error {
	var eg errgroup.Group
	logger.Debug("started accepting DNS queries")

	eg.Go(func() error { return d.dns.ListenAndServe() })
	eg.Go(func() error { return d.dnstls.ListenAndServe() })

	return eg.Wait()
}

// Stop stops the running dns-proxy server
func (d *DNSProxy) Stop() error {
	var eg errgroup.Group
	logger.Debug("stopped accepting DNS queries")

	eg.Go(func() error { return d.dns.Shutdown() })
	eg.Go(func() error { return d.dnstls.Shutdown() })

	return eg.Wait()
}

// NewDNSProxy creates a dns-proxy server
func NewDNSProxy(conf *config.Config) *DNSProxy {
	c := context.Background()

	r := &dnsServer{conf: conf}
	r.Server = &dns.Server{Addr: ":53", Net: "udp", Handler: r}

	t := &dnsServer{conf: conf}
	t.Server = &dns.Server{Addr: ":853", Net: "tcp", Handler: t}
	t.TLSConfig = &tls.Config{GetCertificate: t.GetCertificate}

	d := &DNSProxy{
		conf:   conf,
		acme:   letsencrypt(c),
		dns:    r,
		dnstls: t,
		ctx:    c,
	}
	return d
}
