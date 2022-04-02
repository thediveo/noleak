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
	"strings"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

// IgnoringInBacktrace succeeds if a function name is contained in the backtrace
// of the actual goroutine description.
func IgnoringInBacktrace(fname string) types.GomegaMatcher {
	return &ignoringInBacktraceMatcher{fname: fname}
}

type ignoringInBacktraceMatcher struct {
	fname string
}

// Match succeeds if actual's backtrace contains the specified function name.
func (matcher *ignoringInBacktraceMatcher) Match(actual interface{}) (success bool, err error) {
	g, err := G(actual, "IgnoringInBacktrace")
	if err != nil {
		return false, err
	}
	return strings.Contains(g.Backtrace, matcher.fname), nil
}

// FailureMessage returns a failure message if the actual's backtrace does not
// contain the specified function name.
func (matcher *ignoringInBacktraceMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("to contain %q in the goroutine's backtrace", matcher.fname))
}

// NegatedFailureMessage returns a failure message if the actual's backtrace
// does contain the specified function name.
func (matcher *ignoringInBacktraceMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("not to contain %q in the goroutine's backtrace", matcher.fname))
}
