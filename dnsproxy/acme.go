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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	"github.com/samuelngs/smartdns/config"
	"golang.org/x/crypto/acme"
)

type session struct {
	authz *acme.Authorization
	cha   *acme.Challenge
	label string
	value string
}

type acmeclient struct {
	*acme.Client
	conf *config.Config
}

func letsencrypt(ctx context.Context) *acme.Client {
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Fatal(err.Error())
	}
	c := &acme.Client{
		Key:          k,
		DirectoryURL: "https://acme-v02.api.letsencrypt.org/directory",
	}
	_, err = c.Register(ctx, &acme.Account{}, func(_ string) bool {
		return true
	})
	if err != nil {
		logger.Fatal(err.Error())
	}
	return c
}

func (d *acmeclient) initDNS01Challenge(ctx context.Context) (*session, error) {
	if !d.conf.DNS.TLS.Enabled {
		return nil, errors.New("dns-over-tls is disabled")
	}

	ss := new(session)
	authz, err := d.Authorize(ctx, d.conf.DNS.TLS.Hostname)
	if err != nil {
		return nil, fmt.Errorf("could not authorize acme server: %s", err)
	}
	ss.authz = authz

	for _, c := range authz.Challenges {
		if c.Type == "dns-01" {
			ss.cha = c
		}
	}
	if ss.cha == nil {
		return nil, errors.New("dns-01 challenge is not available")
	}

	ss.label = fmt.Sprintf("_acme-challenge.%s", authz.Identifier.Value)
	ss.value, err = d.DNS01ChallengeRecord(ss.cha.Token)
	if err != nil {
		return nil, fmt.Errorf("could not fetch dns-01 token: %s", err)
	}
	return ss, nil
}

func (d *acmeclient) startDNS01Challenge(ctx context.Context, ss *session) error {
	if !d.conf.DNS.TLS.Enabled {
		return errors.New("dns-over-tls is disabled")
	}
	if ss == nil {
		return errors.New("session is not initialized")
	}
	if _, err := d.Accept(ctx, ss.cha); err != nil {
		return fmt.Errorf("could not accept challenge: %s", err)
	}
	if _, err := d.WaitAuthorization(ctx, ss.authz.URI); err != nil {
		return fmt.Errorf("could not authorize: %s", err)
	}
	return nil
}

func (d *acmeclient) createAcmeCert(ctx context.Context, ss *session) ([][]byte, error) {
	if !d.conf.DNS.TLS.Enabled {
		return nil, errors.New("dns-over-tls is disabled")
	}

	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	r := &x509.CertificateRequest{
		DNSNames: []string{d.conf.DNS.TLS.Hostname},
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, r, k)
	if err != nil {
		return nil, fmt.Errorf("could not initialize certificate request: %s", err)
	}

	crt, _, err := d.CreateCert(ctx, csr, 90*24*time.Hour, true)
	if err != nil {
		return nil, fmt.Errorf("could not create acme certificate: %s", err)
	}

	return crt, nil
}
