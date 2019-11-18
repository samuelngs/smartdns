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

package log

import "strconv"

// FieldType defines the type of field
type FieldType uint8

// Defines the types of fields
const (
	UnknownFieldType FieldType = iota
	BinaryType
	IntegerType
	FloatType
	StringType
	BoolType
)

// Field encapsulates logging fields
type Field struct {
	Key     string    `json:"key"`
	Type    FieldType `json:"type,omitempty"`
	Binary  []byte    `json:"binary,omitempty"`
	Integer int64     `json:"integer,omitempty"`
	Float   float64   `json:"float,omitempty"`
	String  string    `json:"string,omitempty"`
	Bool    bool      `json:"bool,omitempty"`
}

// Value returns field data in string format
func (f Field) Value() string {
	switch f.Type {
	case BinaryType:
		s := string(f.Binary[:])
		return s
	case IntegerType:
		return strconv.FormatInt(f.Integer, 10)
	case FloatType:
		return strconv.FormatFloat(f.Float, 'f', 6, 64)
	case StringType:
		return f.String
	case BoolType:
		if f.Bool {
			return "true"
		}
		return "false"
	default:
		return "unknown"
	}
}

// Tag returns a field that contains string data
func Tag(s string) Field {
	return Field{Type: StringType, String: s}
}

// Binary returns a field that contains bytes data
func Binary(k string, b []byte) Field {
	return Field{Key: k, Type: BinaryType, Binary: b}
}

// String returns a field that contains string data
func String(k string, s string) Field {
	return Field{Key: k, Type: StringType, String: s}
}

// Bool returns a field that contains string data
func Bool(k string, b bool) Field {
	return Field{Key: k, Type: BoolType, Bool: b}
}

// Int returns a field that contains integer data
func Int(k string, i int) Field {
	return Int64(k, int64(i))
}

// Int8 returns a field that contains integer data
func Int8(k string, i int8) Field {
	return Int64(k, int64(i))
}

// Int16 returns a field that contains integer data
func Int16(k string, i int16) Field {
	return Int64(k, int64(i))
}

// Int32 returns a field that contains integer data
func Int32(k string, i int32) Field {
	return Int64(k, int64(i))
}

// Int64 returns a field that contains integer data
func Int64(k string, i int64) Field {
	return Field{Key: k, Type: IntegerType, Integer: i}
}

// Uint returns a field that contains integer data
func Uint(k string, i uint) Field {
	return Uint64(k, uint64(i))
}

// Uint8 returns a field that contains integer data
func Uint8(k string, i uint8) Field {
	return Uint64(k, uint64(i))
}

// Uint16 returns a field that contains integer data
func Uint16(k string, i uint16) Field {
	return Uint64(k, uint64(i))
}

// Uint32 returns a field that contains integer data
func Uint32(k string, i uint32) Field {
	return Uint64(k, uint64(i))
}

// Uint64 returns a field that contains integer data
func Uint64(k string, i uint64) Field {
	return Int64(k, int64(i))
}

// Float32 returns a field that contains float data
func Float32(k string, f float32) Field {
	return Float64(k, float64(f))
}

// Float64 returns a field that contains float data
func Float64(k string, f float64) Field {
	return Field{Key: k, Type: FloatType, Float: f}
}
