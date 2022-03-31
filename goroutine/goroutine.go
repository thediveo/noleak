// Copyright 2022 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package goroutine

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const goroutineHeader = "goroutine "
const goroutineHeaderLen = len(goroutineHeader)

// Goroutine represents information about a single goroutine.
type Goroutine struct {
	ID          int    // goroutine ID
	State       string // goroutine state, such as "running"
	TopFunction string // topmost function on goroutine's stack
	Backtrace   string // goroutine's stack backtrace
}

// String returns a short textual description of this goroutine.
func (g Goroutine) String() string {
	return fmt.Sprintf("Goroutine ID: %d, state: %s, top function: %s",
		g.ID, g.State, g.TopFunction)
}

// Goroutines returns information about all goroutines.
func Goroutines() []Goroutine {
	return goroutines(true)
}

// Current returns information about the current goroutine in which it is
// called.
func Current() Goroutine {
	return goroutines(false)[0]
}

// goroutines is an internal wrapper around dumping either only the stack of the
// current goroutine of the caller or dumping the stacks of all goroutines, and
// then parsing the dump into separate Goroutine descriptions.
func goroutines(all bool) []Goroutine {
	return parseStack(stacks(all))
}

// parseStack parses the stack dump of one or multiple goroutines, as returned
// by runtime.Stack() and then returns a list of Goroutine descriptions based on
// the dump.
func parseStack(stacks []byte) []Goroutine {
	gs := []Goroutine{}

	r := bufio.NewReader(bytes.NewReader(stacks))
	for {
		// We expect a line describing a new "goroutine", everything else is a
		// failure. And yes, if we get an EOF already with this line, bail out.
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		g := new(line)
		// Read the rest ... the backtrace
		g.TopFunction, g.Backtrace = parseGoroutineStack(r)
		gs = append(gs, g)
	}

	return gs
}

// new takes a goroutine line from a stack dump and returns a Goroutine for it.
func new(s string) Goroutine {
	s = strings.TrimSuffix(s, ":\n")
	fields := strings.SplitN(s, " ", 3)
	if len(fields) != 3 {
		panic(fmt.Sprintf("invalid stack header: %q", s))
	}
	id, err := strconv.Atoi(fields[1])
	if err != nil {
		panic(fmt.Sprintf("invalid stack header ID: %q, header: %q", fields[1], s))
	}
	state := strings.TrimSuffix(strings.TrimPrefix(fields[2], "["), "]")
	return Goroutine{ID: id, State: state}
}

// parseGoroutineStack reads from stack information from a reader until the next
// goroutine header is seen. The next goroutine header isn't consumed so that
// the caller can still read the next header.
func parseGoroutineStack(r *bufio.Reader) (topF string, backtrace string) {
	stack := bytes.Buffer{}
	// Read stack information belonging to this goroutine until we meet
	// another goroutine header.
	for {
		header, err := r.Peek(goroutineHeaderLen)
		if string(header) == goroutineHeader {
			break
		}
		if err != nil && err != io.EOF {
			panic("parsing stack backtrace failed: " + err.Error())
		}
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			panic("parsing stack backtrace failed: " + err.Error())
		}
		// The first line after a goroutine header lists the "topmost" function.
		if topF == "" {
			line := /*sic!*/ strings.TrimSpace(line)
			idx := strings.LastIndex(line, "(")
			if idx <= 0 {
				panic(fmt.Sprintf("invalid function call stack entry: %q", line))
			}
			topF = line[:idx]
		}
		// Always append the line to the goroutine's stack backtrace.
		stack.WriteString(line)
		if err == io.EOF {
			break
		}
	}
	return topF, stack.String()
}
