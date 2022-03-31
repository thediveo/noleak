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

package noleak

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/onsi/gomega/format"
	"github.com/thediveo/noleak/goroutine"
)

// G takes an actual "any" untyped value and returns it as a typed Goroutine, if
// possible. It returns an error if actual isn't of either type Goroutine or a
// pointer to it. G is intended to be mainly used by goroutine-related Gomega
// matchers, such as IgnoringTopFunction, et cetera.
func G(actual interface{}, matchername string) (goroutine.Goroutine, error) {
	if actual != nil {
		switch actual := actual.(type) {
		case goroutine.Goroutine:
			return actual, nil
		case *goroutine.Goroutine:
			return *actual, nil
		}
	}
	return goroutine.Goroutine{},
		fmt.Errorf("%s matcher expects a goroutine.Goroutine or *goroutine.Goroutine.  Got:\n%s",
			matchername, format.Object(actual, 1))
}

// goids returns a (sorted) list of Goroutine IDs in textual format.
func goids(gs []goroutine.Goroutine) string {
	ids := make([]int, len(gs))
	for idx, g := range gs {
		ids[idx] = g.ID
	}
	sort.Ints(ids)
	var buff strings.Builder
	for idx, id := range ids {
		if idx > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(strconv.FormatInt(int64(id), 10))
	}
	return buff.String()
}
