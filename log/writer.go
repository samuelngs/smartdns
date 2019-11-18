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

import (
	"strings"
)

// Composer encapsulates a text writer
type Composer interface {
	Write(string)
	Bold(string) string
	Underline(string) string
	Blink(string) string
	Grey(string) string
	Red(string) string
	Green(string) string
	Yellow(string) string
	Blue(string) string
	LightGrey(string) string
	String() string
}

type composer struct {
	sb strings.Builder
}

func (c *composer) Write(s string) {
	c.sb.WriteString(s)
	c.sb.WriteString("\033[0m")
}

func (c *composer) styled(a, s string) string {
	return a + s
}

func (c *composer) Bold(s string) string {
	return c.styled("\033[1m", s)
}

func (c *composer) Underline(s string) string {
	return c.styled("\033[4m", s)
}

func (c *composer) Blink(s string) string {
	return c.styled("\033[5m", s)
}

func (c *composer) Grey(s string) string {
	return c.styled("\x1b[30m", s)
}

func (c *composer) Red(s string) string {
	return c.styled("\x1b[91m", s)
}

func (c *composer) Green(s string) string {
	return c.styled("\x1b[92m", s)
}

func (c *composer) Yellow(s string) string {
	return c.styled("\x1b[93m", s)
}

func (c *composer) Blue(s string) string {
	return c.styled("\x1b[94m", s)
}

func (c *composer) LightGrey(s string) string {
	return c.styled("\x1b[37m", s)
}

func (c *composer) String() string {
	return c.sb.String() + "\n"
}

// NewWriter returns a new text writer
func NewWriter() Composer {
	var sb strings.Builder
	return &composer{
		sb: sb,
	}
}
