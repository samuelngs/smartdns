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

package https

import (
	"bytes"
	"io"
	"net"
)

// Handshake on a https request
type Handshake struct {
	Hostname string
	Buffer   bytes.Buffer
}

// ParseHandshakeMessage for parsing handshake metadata on a https request
func ParseHandshakeMessage(c *net.TCPConn) (*Handshake, error) {
	var buf bytes.Buffer
	buf.WriteByte(22)

	r := io.MultiReader(bytes.NewReader([]byte{22}), io.TeeReader(c, &buf))
	hr := &handshakeReader{r: r}

	hostname, err := hr.ReadSNIHostname()
	if err != nil {
		return nil, err
	}

	h := &Handshake{
		Hostname: hostname,
		Buffer:   buf,
	}
	return h, nil
}
