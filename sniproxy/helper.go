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

	"golang.org/x/sync/errgroup"
)

type reader struct {
	io.Reader
	ResetTimeout func()
	Close        func() error
	RemoteAddr   func() net.Addr
}

type writer struct {
	io.Writer
	ResetTimeout func()
	Close        func() error
	RemoteAddr   func() net.Addr
}

func newReader(c *net.TCPConn, timeout time.Duration, prefix ...io.Reader) reader {
	var rd io.Reader
	if len(prefix) > 0 && prefix[0] != nil {
		rd = io.MultiReader(prefix[0], c)
	} else {
		rd = c
	}
	return reader{
		Reader:       rd,
		ResetTimeout: func() { c.SetReadDeadline(time.Now().Add(timeout)) },
		Close:        c.CloseRead,
		RemoteAddr:   c.RemoteAddr,
	}
}

func newWriter(c *net.TCPConn, timeout time.Duration) writer {
	return writer{
		Writer:       c,
		ResetTimeout: func() { c.SetWriteDeadline(time.Now().Add(timeout)) },
		Close:        c.CloseWrite,
		RemoteAddr:   c.RemoteAddr,
	}
}

func proxy(src, dst *net.TCPConn, timeout time.Duration, prefix io.Reader) error {
	var eg errgroup.Group
	eg.Go(func() error { return forward(newWriter(dst, timeout), newReader(src, timeout, prefix)) })
	eg.Go(func() error { return forward(newWriter(src, timeout), newReader(dst, timeout)) })
	return eg.Wait()
}

func forward(dst writer, src reader) error {
	defer src.Close()
	defer dst.Close()

	var buf [4096]byte
	defer src.ResetTimeout()
	defer dst.ResetTimeout()

	for {
		n, err := src.Read(buf[:])
		if n > 0 {
			src.ResetTimeout()
			if _, err := dst.Write(buf[:n]); err != nil {
				return fmt.Errorf("could not write to %q: %v", dst.RemoteAddr(), err)
			}
			dst.ResetTimeout()
		}
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return fmt.Errorf("error reading from %q: %v", src.RemoteAddr(), err)
		}
	}
}
