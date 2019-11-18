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
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

func readUint8(r io.Reader) (uint8, error) {
	var buf [1]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, err
	}
	return buf[0], nil
}

// readUint16 interprets data as big-endian.
func readUint16(r io.Reader) (uint16, error) {
	var buf [2]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(buf[:]), nil
}

// readUint24 interprets data as big-endian.
func readUint24(r io.Reader) (uint32, error) {
	var buf [3]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, err
	}
	return uint32(buf[0])<<16 | uint32(buf[1])<<8 | uint32(buf[2]), nil
}

func skip(r io.Reader, sz int64) error {
	_, err := io.CopyN(ioutil.Discard, r, sz)
	return err
}

func skipVec8(r io.Reader) error {
	vl, err := readUint8(r)
	if err != nil {
		return fmt.Errorf("could not read length: %v", err)
	}
	if err := skip(r, int64(vl)); err != nil {
		return fmt.Errorf("could not skip content: %v", err)
	}
	return nil
}

func skipVec16(r io.Reader) error {
	vl, err := readUint16(r)
	if err != nil {
		return fmt.Errorf("could not read length: %v", err)
	}
	if err := skip(r, int64(vl)); err != nil {
		return fmt.Errorf("could not skip content: %v", err)
	}
	return nil
}
