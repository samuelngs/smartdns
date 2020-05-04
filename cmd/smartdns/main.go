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

package main

import (
	"github.com/samuelngs/smartdns/config"
	"github.com/samuelngs/smartdns/dnsproxy"
	"github.com/samuelngs/smartdns/log"
	"github.com/samuelngs/smartdns/sniproxy"
	"golang.org/x/sync/errgroup"
)

var logger = log.DefaultLogger

func main() {
	var eg errgroup.Group

	conf := config.DefaultConfig()
	conf.DNS.TLS.Enabled = false

	sniproxy := sniproxy.NewSNIProxy(conf)
	dnsproxy := dnsproxy.NewDNSProxy(conf)

	eg.Go(func() error { return sniproxy.Start() })
	eg.Go(func() error { return dnsproxy.Start() })

	if err := eg.Wait(); err != nil {
		logger.Fatal(err.Error())
	}
}
