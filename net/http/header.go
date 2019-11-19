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

package http

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
)

var hostHeaderPrefix = []byte("Host:")

// ParseHost extracts host from a http request header
func ParseHost(c *net.TCPConn) (string, io.Reader, error) {
	var buf bytes.Buffer
	var hostname string

	sc := bufio.NewScanner(io.TeeReader(c, &buf))
	sc.Scan()
lr:
	for sc.Scan() {
		switch ln := sc.Bytes(); {
		case len(ln) == 0:
			break lr
		case bytes.HasPrefix(ln, hostHeaderPrefix):
			hostname = string(bytes.TrimSpace(bytes.TrimPrefix(ln, hostHeaderPrefix)))
			break lr
		}
	}
	if err := sc.Err(); err != nil {
		return "", nil, errors.New("could not read request body")
	}
	if len(hostname) == 0 {
		return "", nil, errors.New("could not read host header")
	}

	return hostname, &buf, nil
}
