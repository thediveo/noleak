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
	"sort"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"github.com/thediveo/noleak/goroutine"
)

// IgnoringGoroutines succeeds if an actual goroutine, identified by its ID, is
// in a slice of expected goroutines. A typical use of the IgnoringGoroutines
// matcher is to take a snapshot of the current goroutines just right before a
// test and then at the end of a test filtering out these "good" and known
// goroutines.
func IgnoringGoroutines(goroutines []goroutine.Goroutine) types.GomegaMatcher {
	m := &ignoringGoroutinesMatcher{
		ignoreGoids: map[int]struct{}{},
	}
	for _, g := range goroutines {
		m.ignoreGoids[g.ID] = struct{}{}
	}
	return m
}

type ignoringGoroutinesMatcher struct {
	ignoreGoids map[int]struct{}
}

// Match succeeds if actual is a goroutine.Goroutine and its ID is in the set of
// goroutine IDs to expect and thus to ignore in leak checks.
func (matcher *ignoringGoroutinesMatcher) Match(actual interface{}) (success bool, err error) {
	g, err := G(actual, "IgnoringGoroutines")
	if err != nil {
		return false, err
	}
	_, ok := matcher.ignoreGoids[g.ID]
	return ok, nil
}

// FailureMessage returns a failure message if the actual goroutine isn't in the
// set of goroutines to be ignored.
func (matcher *ignoringGoroutinesMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to be contained in the list of expected goroutine IDs", matcher.expectedGoids())
}

// NegatedFailureMessage returns a negated failure message if the actual
// goroutine actually is in the set of goroutines to be ignored.
func (matcher *ignoringGoroutinesMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to be contained in the list of expected goroutine IDs", matcher.expectedGoids())
}

// expectedGoids returns the sorted list of expected goroutine IDs.
func (matcher *ignoringGoroutinesMatcher) expectedGoids() []int {
	ids := make([]int, 0, len(matcher.ignoreGoids))
	for id := range matcher.ignoreGoids {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	return ids
}
