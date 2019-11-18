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
	"fmt"
	"os"
	"strings"
)

// lvl sets current logging level
var lvl = LogInfo

// Level defines the importance and urgency of the log message
type Level int

// Defines the importance level of logs
const (
	LogFatal Level = iota
	LogError
	LogWarn
	LogInfo
	LogDebug
	LogTrace
)

var logLevelRefs = map[Level]string{
	LogFatal: "fatal",
	LogError: "error",
	LogWarn:  "warn",
	LogInfo:  "info",
	LogDebug: "debug",
	LogTrace: "trace",
}

var logLevelIds = map[string]Level{
	"fatal": LogFatal,
	"error": LogError,
	"warn":  LogWarn,
	"info":  LogInfo,
	"debug": LogDebug,
	"trace": LogTrace,
}

// Head pads log level reference to 5 characters, this would
// left-justifies the strings and add spaces to fill the empty
// columns on the right.
func (l Level) Head() string {
	if s, ok := logLevelRefs[l]; ok {
		return fmt.Sprintf("%-5s", strings.ToUpper(s))
	}
	panic("unrecognized log level")
}

func (l Level) String() string {
	if s, ok := logLevelRefs[l]; ok {
		return s
	}
	panic("unrecognized log level")
}

// LevelFromString returns the log level enum from a string
func LevelFromString(s string) Level {
	if l, ok := logLevelIds[strings.ToLower(s)]; ok {
		return l
	}
	panic("unrecognized log level")
}

// SetLevel sets the logging verbose level
func SetLevel(l Level) {
	lvl = l
}

func init() {
	l, ok := os.LookupEnv("LOG")
	switch {
	case l == "*":
		lvl = LogTrace
	case ok:
		lvl = LevelFromString(l)
	default:
		lvl = LogInfo
	}
}
