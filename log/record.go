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
	"time"

	"github.com/go-stack/stack"
)

const skipLevel = 2

// Record encapsulates a log event
type Record struct {
	Time   time.Time  `json:"time"`
	Level  Level      `json:"level"`
	Msg    string     `json:"message"`
	Fields []Field    `json:"fields"`
	Call   stack.Call `json:"call"`
}

func (r *Record) String() string {
	c := NewWriter()
	c.Write(c.Green(r.Level.Head()))
	c.Write(c.Green(" | "))
	c.Write(c.Grey(c.Bold(r.Time.Format("2006-01-02T15:04:05"))))
	c.Write(" ")
	c.Write(r.Msg)
	c.Write(" ")
	for _, field := range r.Fields {
		c.Write(c.LightGrey("["))
		if len(field.Key) > 0 {
			c.Write(c.LightGrey(field.Key))
			c.Write(c.LightGrey(": "))
		}
		c.Write(c.LightGrey(field.Value()))
		c.Write(c.LightGrey("] "))
	}
	return c.String()
}
