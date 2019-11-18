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
	"time"

	"github.com/go-stack/stack"
	"github.com/samuelngs/smartdns/log"
	"github.com/stretchr/testify/assert"
)

func TestRecordOutput(t *testing.T) {
	record := &log.Record{
		Time:  time.Date(2019, time.September, 21, 0, 0, 0, 0, time.UTC),
		Level: log.LogInfo,
		Msg:   "hello world",
		Fields: []log.Field{
			log.String("key", "val"),
		},
		Call: stack.Caller(2),
	}
	assert.Equal(t,
		"\x1b[92mINFO \x1b[0m\x1b[92m | \x1b[0m\x1b[30m\x1b[1m2019-09-21T00:00:00\x1b[0m \x1b[0mhello world\x1b[0m \x1b[0m\x1b[37m[\x1b[0m\x1b[37mkey\x1b[0m\x1b[37m: \x1b[0m\x1b[37mval\x1b[0m\x1b[37m] \x1b[0m\n",
		record.String())
}
