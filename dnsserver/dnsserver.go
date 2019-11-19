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

package dnsserver

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/samuelngs/smartdns/config"
	"github.com/samuelngs/smartdns/log"
)

var logger = log.DefaultLogger

type dnsServer struct {
	conf *config.Config
}

func (d *dnsServer) parseQuery(r *dns.Msg) (dns.Question, bool) {
	if r.Opcode == dns.OpcodeQuery {
		for _, q := range r.Question {
			if q.Qtype == dns.TypeA {
				return q, true
			}
		}
	}
	return dns.Question{}, false
}

func (d *dnsServer) handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.Compress = false
	m.SetReply(r)
	defer w.WriteMsg(m)
	defer logger.Trace("dns query completed")

	if !d.conf.Network.IsAllowedIP(w.RemoteAddr()) {
		return
	}
	question, ok := d.parseQuery(r)
	if !ok {
		return
	}

	logger.Trace(
		"dns query accepted",
		log.String("name", question.Name))

	list := config.DNSResolveList(d.conf.Proxy.DNS)
	resolv := list.MatchDNS(question.Name)

	var ttl = 60
	if resolv != nil && resolv.TTL != ttl && resolv.TTL > 0 {
		ttl = resolv.TTL
	}

	switch {
	case resolv != nil && len(resolv.Nameserver) > 0:
		logger.Trace(
			"resolving domain name with nameserver",
			log.String("name", question.Name),
			log.String("nameserver", resolv.Nameserver))

		t := new(dns.Msg)
		t.SetQuestion(question.Name, dns.TypeA)
		c := new(dns.Client)
		if in, _, _ := c.Exchange(t, resolv.NameserverAddr()); len(in.Answer) > 0 {
			for _, a := range in.Answer {
				r, _ := dns.NewRR(a.String())
				m.Answer = append(m.Answer, r)
			}
		}

	case resolv != nil && len(resolv.IP) > 0:
		logger.Trace(
			"resolving domain name to ip",
			log.String("name", question.Name),
			log.String("ip", resolv.IP))

		r, _ := dns.NewRR(fmt.Sprintf("%s %d IN A %s", question.Name, ttl, resolv.IP))
		m.Answer = []dns.RR{r}

	default:
		logger.Trace(
			"resolving domain name to proxy ip",
			log.String("name", question.Name),
			log.String("ip", d.conf.Proxy.Host))

		r, _ := dns.NewRR(fmt.Sprintf("%s %d IN A %s", question.Name, ttl, d.conf.Proxy.Host))
		m.Answer = []dns.RR{r}
	}
}

// Listen starts dns server
func Listen(conf *config.Config) {
	d := &dnsServer{conf}
	dns.HandleFunc(".", d.handleDNSRequest)
	server := &dns.Server{Addr: ":53", Net: "udp"}

	logger.Debug("accepting DNS queries")
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal(err.Error())
	}
}
