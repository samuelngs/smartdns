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
	"fmt"
	"sync"

	"github.com/miekg/dns"
	"github.com/samuelngs/smartdns/config"
	"github.com/samuelngs/smartdns/log"
)

type dnsServer struct {
	*dns.Server
	txt  *sync.Map
	conf *config.Config
}

func (d *dnsServer) parseQuery(r *dns.Msg) (dns.Question, bool) {
	if r.Opcode == dns.OpcodeQuery {
		for _, q := range r.Question {
			if q.Qtype == dns.TypeA {
				return q, true
			}
			if q.Qtype == dns.TypeTXT {
				return q, true
			}
		}
	}
	return dns.Question{}, false
}

func (d *dnsServer) resolveA(m *dns.Msg, question dns.Question) {
	list := config.DNSResolveList(d.conf.DNS.DNSResolveList)
	resolv := list.MatchDNS(question.Name)

	var ttl = 60
	if resolv != nil && resolv.TTL != ttl && resolv.TTL > 0 {
		ttl = resolv.TTL
	}

	switch {
	case resolv != nil && resolv.Nameserver == "-":
		logger.Trace(
			"resolving domain name to proxy ip",
			log.String("name", question.Name),
			log.String("ip", d.conf.SNIProxy.Host))

		r, _ := dns.NewRR(fmt.Sprintf("%s %d IN A %s", question.Name, ttl, d.conf.SNIProxy.Host))
		m.Answer = []dns.RR{r}

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
			"resolving domain name with google nameserver",
			log.String("name", question.Name))

		t := new(dns.Msg)
		t.SetQuestion(question.Name, dns.TypeA)
		c := new(dns.Client)
		if in, _, _ := c.Exchange(t, "8.8.8.8:53"); len(in.Answer) > 0 {
			for _, a := range in.Answer {
				r, _ := dns.NewRR(a.String())
				m.Answer = append(m.Answer, r)
			}
		}
	}
}

func (d *dnsServer) resolveTXT(m *dns.Msg, question dns.Question) {
	if o, ok := d.txt.Load(question.Name); ok {
		r, _ := dns.NewRR(fmt.Sprintf(`%s %d IN TXT "%s"`, question.Name, 60, o.(string)))
		m.Answer = []dns.RR{r}
	}
}

func (d *dnsServer) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
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

	switch question.Qtype {
	case dns.TypeA:
		d.resolveA(m, question)
	case dns.TypeTXT:
		d.resolveTXT(m, question)
	}
}
