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
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	maxRecordSize        = 1 << 14
	contentTypeHandshake = 22
	handshakeTypeHello   = 1
	extensionTypeSNI     = 0
	nameTypeHostName     = 0
)

type handshakeReader struct {
	r io.Reader
	l int
	e error
}

func (r *handshakeReader) RecordError(err error) {
	if err != nil {
		r.e = err
	}
}

func (r *handshakeReader) Read(b []byte) (length int, err error) {
	defer r.RecordError(err)
	if r.e != nil {
		return 0, r.e
	}

	if r.l == 0 {
		typ, err := readUint8(r.r)
		if err != nil {
			return 0, fmt.Errorf("could not read record header: %v", err)
		}
		if typ != contentTypeHandshake {
			return 0, fmt.Errorf("got wrong content type (wanted %d, got %d)", contentTypeHandshake, typ)
		}
		if err := skip(r.r, 2); err != nil {
			return 0, fmt.Errorf("could not read record header: %v", err)
		}
		sz, err := readUint16(r.r)
		if err != nil {
			return 0, fmt.Errorf("could not read record header: %v", err)
		}
		if sz > maxRecordSize {
			return 0, fmt.Errorf("record too large (%d > %d)", sz, maxRecordSize)
		}
		r.l = int(sz)
	}

	if len(b) > r.l {
		b = b[:r.l]
	}
	n, err := r.r.Read(b)
	r.l -= n
	return n, err
}

func (r *handshakeReader) ReadExtensions() (io.Reader, error) {
	var rd io.Reader = r

	typ, err := readUint8(rd)
	if err != nil {
		return nil, fmt.Errorf("could not read msg_type: %v", err)
	}
	if typ != handshakeTypeHello {
		return nil, fmt.Errorf("handshake message not a ClientHello (type %d, expected %d)", typ, handshakeTypeHello)
	}

	l, err := readUint24(rd)
	if err != nil {
		return nil, fmt.Errorf("could not read handshake message length: %v", err)
	}
	rd = io.LimitReader(rd, int64(l))

	// skip the protocol version (2 bytes)
	if err := skip(rd, 2); err != nil {
		return nil, fmt.Errorf("could not skip client_version: %v", err)
	}

	// skip the clienthello.random (32 bytes, out of which 28 are suppose
	// to be generated with a cryptographically strong number generator)
	if err := skip(rd, 32); err != nil {
		return nil, fmt.Errorf("could not skip random: %v", err)
	}

	// the "session_id" (in case the client wants to resume a session in
	// an abbreviated handshake, see below)
	if err := skipVec8(rd); err != nil {
		return nil, fmt.Errorf("could not skip session_id: %v", err)
	}

	// skip the list of "cipher suites" that the client knows of, ordered
	// by client preference
	if err := skipVec16(rd); err != nil {
		return nil, fmt.Errorf("could not skip cipher_suites: %v", err)
	}

	// the list of compression algorithms that the client knows of, ordered
	// by client preference
	if err := skipVec8(rd); err != nil {
		return nil, fmt.Errorf("could not skip compression_methods: %v", err)
	}

	// read extensions
	n, err := readUint16(rd)
	if err != nil {
		return nil, errors.New("no extension is available")
	}
	if err != nil {
		return nil, fmt.Errorf("could not read extensions length: %v", err)
	}

	rd = io.LimitReader(rd, int64(n))
	return rd, nil
}

func (r *handshakeReader) ReadSNIHostname() (string, error) {
	rd, err := r.ReadExtensions()
	if err != nil {
		return "", err
	}
	for {
		typ, err := readUint16(rd)
		if err == io.EOF {
			return "", errors.New("no SNI extension")
		}
		if err != nil {
			return "", fmt.Errorf("could not read extension_type: %v", err)
		}
		n, err := readUint16(rd)
		if err != nil {
			return "", fmt.Errorf("could not read extension_data length: %v", err)
		}
		if typ != extensionTypeSNI {
			if err := skip(rd, int64(n)); err != nil {
				return "", fmt.Errorf("could not skip extension_data: %v", err)
			}
			continue
		}

		er := io.LimitReader(rd, int64(n))
		sl, err := readUint16(er)
		if err != nil {
			return "", fmt.Errorf("could not read server_name_list length: %v", err)
		}

		er = io.LimitReader(rd, int64(sl))
		for {
			typ, err := readUint8(er)
			if err == io.EOF {
				return "", errors.New("SNI extension has no ServerName of type host_name")
			}
			if err != nil {
				return "", fmt.Errorf("could not read name_type: %v", err)
			}
			if typ != nameTypeHostName {
				if err := skipVec16(rd); err != nil {
					return "", fmt.Errorf("could not skip server_name_list entry: %v", err)
				}
				continue
			}

			nl, err := readUint16(er)
			if err != nil {
				return "", fmt.Errorf("could not read host_name length: %v", err)
			}
			var b strings.Builder
			if _, err := io.CopyN(&b, er, int64(nl)); err != nil {
				return "", fmt.Errorf("could not read HostName: %v", err)
			}
			return b.String(), nil
		}
	}
}
