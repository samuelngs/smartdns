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
	"sync"

	"github.com/miekg/dns"
	"github.com/samuelngs/smartdns/config"
	"github.com/samuelngs/smartdns/log"
	"golang.org/x/sync/errgroup"
)

var logger = log.DefaultLogger

// DNSProxy constructs a dns-proxy server
type DNSProxy struct {
	conf   *config.Config
	acme   *acmeclient
	dns    *dnsServer
	dnstls *dnsServer
	ctx    context.Context
}

// Start initializes and starts dns-proxy server
func (d *DNSProxy) Start() error {
	var eg errgroup.Group
	logger.Debug("started accepting DNS queries")

	eg.Go(func() error { return d.dns.ListenAndServe() })
	eg.Go(d.startDOTServer)

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

func (d *DNSProxy) startDOTServer() error {
	logger.Debug("initialize dns-01 challenge")
	s, err := d.acme.initDNS01Challenge(d.ctx)
	if err != nil {
		return err
	}

	logger.Debug(
		"set up dns-01 challenge verification",
		log.String("label", s.label),
		log.String("value", s.value))
	d.dns.txt.Store(s.label, s.value)
	defer d.dns.txt.Delete(s.label)

	logger.Debug("start dns-01 verification")
	if err := d.acme.startDNS01Challenge(d.ctx, s); err != nil {
		return err
	}

	logger.Debug("create acme certificate")
	o, err := d.acme.createAcmeCert(d.ctx, s)
	if err != nil {
		return err
	}
	if o == nil {
		return nil
	}

	d.dnstls.Shutdown()
	d.dnstls.Server.TLSConfig = &tls.Config{}
	return d.dnstls.ListenAndServe()
}

// NewDNSProxy creates a dns-proxy server
func NewDNSProxy(conf *config.Config) *DNSProxy {
	m := new(sync.Map)
	c := context.Background()

	a := letsencrypt(c)
	a.withConfig(conf)

	r := &dnsServer{conf: conf, txt: m}
	r.Server = &dns.Server{Addr: ":53", Net: "udp", Handler: r}

	t := &dnsServer{conf: conf, txt: m}
	t.Server = &dns.Server{Addr: ":853", Net: "tcp", Handler: t}

	return &DNSProxy{
		conf:   conf,
		acme:   a,
		dns:    r,
		dnstls: t,
		ctx:    c,
	}
}
