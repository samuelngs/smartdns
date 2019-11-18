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

package log_test

import (
	"testing"

	"github.com/samuelngs/smartdns/log"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	writer := log.NewWriter()
	writer.Write("test")
	assert.Equal(t, "test\033[0m\n", writer.String())
}

func TestStyles(t *testing.T) {
	writer := log.NewWriter()
	assert.Equal(t, "\033[1mtest", writer.Bold("test"))
	assert.Equal(t, "\033[4mtest", writer.Underline("test"))
	assert.Equal(t, "\033[5mtest", writer.Blink("test"))
}

func TestColors(t *testing.T) {
	writer := log.NewWriter()
	assert.Equal(t, "\x1b[30mtest", writer.Grey("test"))
	assert.Equal(t, "\x1b[91mtest", writer.Red("test"))
	assert.Equal(t, "\x1b[92mtest", writer.Green("test"))
	assert.Equal(t, "\x1b[93mtest", writer.Yellow("test"))
	assert.Equal(t, "\x1b[94mtest", writer.Blue("test"))
}

func TestChainStyles(t *testing.T) {
	writer := log.NewWriter()
	styled := writer.Bold(writer.Red("test"))
	assert.Equal(t, "\033[1m\x1b[91mtest", styled)
}

func TestChainWrite(t *testing.T) {
	writer := log.NewWriter()
	writer.Write(writer.Bold(writer.Red("1")))
	writer.Write(writer.Underline(writer.Red("2")))
	assert.Equal(t, "\x1b[1m\x1b[91m1\x1b[0m\x1b[4m\x1b[91m2\x1b[0m\n", writer.String())
}
