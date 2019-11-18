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

func TestLevelHead(t *testing.T) {
	assert.Equal(t, "FATAL", log.LogFatal.Head())
	assert.Equal(t, "ERROR", log.LogError.Head())
	assert.Equal(t, "WARN ", log.LogWarn.Head())
	assert.Equal(t, "INFO ", log.LogInfo.Head())
	assert.Equal(t, "DEBUG", log.LogDebug.Head())
	assert.Equal(t, "TRACE", log.LogTrace.Head())
	assert.Panics(t, func() {
		_ = log.Level(9999).Head()
	})
}

func TestLevelString(t *testing.T) {
	assert.Equal(t, "fatal", log.LogFatal.String())
	assert.Equal(t, "error", log.LogError.String())
	assert.Equal(t, "warn", log.LogWarn.String())
	assert.Equal(t, "info", log.LogInfo.String())
	assert.Equal(t, "debug", log.LogDebug.String())
	assert.Equal(t, "trace", log.LogTrace.String())
	assert.Panics(t, func() {
		_ = log.Level(9999).String()
	})
}
