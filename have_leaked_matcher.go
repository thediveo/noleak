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
	"reflect"
	"strings"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"github.com/thediveo/noleak/goroutine"
)

// standardFilters specifies the always automatically included no-leak goroutine
// filter matchers.
//
// Note: it's okay to instantiate the Gomega Matchers here, as all goroutine
// filtering-related noleak matchers are stateless with respect to any actual
// value they try to match. This allows us to simply prepend them to any
// user-supplied optional matchers when HaveLeaked returns a new goroutine
// leakage detecting matcher.
//
// Note: cgo's goroutines with status "[syscall, locked to thread]" do not
// appear any longer (since mid-2017), as these cgo goroutines are put into the
// "dead" state when not in use. See: https://github.com/golang/go/issues/16714
// and https://go-review.googlesource.com/c/go/+/45030/.
var standardFilters = []types.GomegaMatcher{
	// Ginkgo testing framework
	IgnoringTopFunction("github.com/onsi/ginkgo/v2/internal.(*Suite).runNode"),
	IgnoringTopFunction("github.com/onsi/ginkgo/v2/internal.(*Suite).runNode..."),
	IgnoringTopFunction("github.com/onsi/ginkgo/v2/internal/interrupt_handler.(*InterruptHandler).registerForInterrupts..."),
	IgnoringTopFunction("github.com/onsi/ginkgo/internal/specrunner.(*SpecRunner).registerForInterrupts"),

	// goroutines of Go's own testing package for its own workings...
	IgnoringTopFunction("testing.RunTests [chan receive]"),
	IgnoringTopFunction("testing.(*T).Run [chan receive]"),
	IgnoringTopFunction("testing.(*T).Parallel [chan receive]"),

	// os/signal starts its own runtime goroutine, where loop calls signal_recv
	// in a loop, so we need to expect them both...
	IgnoringTopFunction("os/signal.signal_recv"),
	IgnoringTopFunction("os/signal.loop"),

	// signal.Notify starts a runtime goroutine...
	IgnoringInBacktrace("runtime.ensureSigM"),

	// reading a trace...
	IgnoringInBacktrace("runtime.ReadTrace"),
}

// HaveLeaked succeeds (or rather, "suckceeds" considering it appears in failing
// tests) if after filtering out ("ignoring") the expected goroutines from the
// list of actual goroutines the remaining list of goroutines is non-empty.
// These goroutines not filtered out are considered to have been leaked.
//
// For convenience, HaveLeaked automatically filters out well-known runtime and
// testing goroutines using a built-in standard filter matchers list. In
// addition to the built-in filters, HaveLeaked accepts an optional list of
// non-leaky goroutine filter matchers. These filtering matchers can be
// specified in different formats, as described below.
//
// Since there might be "pending" goroutines at the end of tests that eventually
// will properly wind down so they aren't leaking, HaveLeaked is best paired
// with Eventually instead of Expect. In its shortest form this will use
// Eventually's default timeout and polling interval settings, but these can be
// overridden as usual:
//
//   // Remember to use "Goroutines" and not "Goroutines()" with Eventually()!
//   Eventually(Goroutines).ShouldNot(HaveLeaked())
//   Eventually(Goroutines).WithTimeout(5 * time.Second).ShouldNot(HaveLeaked())
//
// In its simplest form, an expected non-leaky goroutine can be identified by
// passing the (fully qualified) name (in form of a string) of the topmost
// function on the backtrace stack. For instance:
//
//   Eventually(Goroutines).ShouldNot(HaveLeaked("foo.bar"))
//
// This is the shorthand equivalent to this explicit form:
//
//   Eventually(Goroutines).ShouldNot(HaveLeaked(IgnoringTopFunction("foo.bar")))
//
// HaveLeak also accepts passing a slice of Goroutine objects to be considered
// non-leaky goroutines.
//
//   snapshot := Goroutines()
//   DoSomething()
//   Eventually(Goroutines).ShouldNot(HaveLeaked(snapshot))
//
// Again, this is shorthand for the following explicit form:
//
//   snapshot := Goroutines()
//   DoSomething()
//   Eventually(Goroutines).ShouldNot(HaveLeaked(IgnoringGoroutines(snapshot)))
//
// Finally, HaveLeaked accepts any GomegaMatcher and will repeatedly pass it a
// Goroutine object: if the matcher succeeds, the Goroutine object in question
// is considered to be non-leaked and thus filtered out. While the following
// built-in Goroutine filter matchers should hopefully cover most situations,
// any suitable GomegaMatcher can be used for tricky leaky Goroutine filtering.
//
//   IgnoringTopFunction("foo.bar")
//   IgnoringTopFunction("foo.bar...")
//   IgnoringTopFunction("foo.bar [chan receive]")
//   IgnoringGoroutines(expectedGoroutines)
//   IgnoringInBacktrace("foo.bar.baz")
func HaveLeaked(ignoring ...interface{}) types.GomegaMatcher {
	m := &HaveLeakedMatcher{filters: standardFilters}
	for _, ign := range ignoring {
		switch ign := ign.(type) {
		case string:
			m.filters = append(m.filters, IgnoringTopFunction(ign))
		case []goroutine.Goroutine:
			m.filters = append(m.filters, IgnoringGoroutines(ign))
		case types.GomegaMatcher:
			m.filters = append(m.filters, ign)
		default:
			panic(fmt.Sprintf("HaveLeaked expected a string, []Goroutine, or GomegaMatcher, but got:\n%s", format.Object(ign, 1)))
		}
	}
	return m
}

// HaveLeakedMatcher implements the HaveLeaked Gomega Matcher that succeeds if
// the actual list of goroutines is non-empty after filtering out the expected
// goroutines.
type HaveLeakedMatcher struct {
	filters []types.GomegaMatcher // expected goroutines that aren't leaks.
	leaked  []goroutine.Goroutine // surplus goroutines which we consider to be leaks.
}

var gsT = reflect.TypeOf([]goroutine.Goroutine{})

// Match succeeds if actual is an array or slice of goroutine.Goroutine
// information and still contains goroutines after filtering out all expected
// goroutines that were specified when creating the matcher.
func (matcher *HaveLeakedMatcher) Match(actual interface{}) (success bool, err error) {
	val := reflect.ValueOf(actual)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		if !val.Type().AssignableTo(gsT) {
			return false, fmt.Errorf(
				"HaveLeaked matcher expects an array or slice of goroutines.  Got:\n%s",
				format.Object(actual, 1))
		}
	default:
		return false, fmt.Errorf(
			"HaveLeaked matcher expects an array or slice of goroutines.  Got:\n%s",
			format.Object(actual, 1))
	}
	goroutines := val.Convert(gsT).Interface().([]goroutine.Goroutine)
	matcher.leaked, err = matcher.filter(goroutines, matcher.filters)
	if err != nil {
		return false, err
	}
	if len(matcher.leaked) == 0 {
		return false, nil
	}
	return true, nil // we have leak(ed)
}

// FailureMessage returns a failure message if there are leaked goroutines.
func (matcher *HaveLeakedMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected to leak goroutines:\n%s", matcher.listGoroutines(matcher.leaked, 1))
}

// NegatedFailureMessage returns a negated failure message if there aren't any leaked goroutines.
func (matcher *HaveLeakedMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected not to leak goroutines:\n%s", matcher.listGoroutines(matcher.leaked, 1))
}

// listGoroutines returns a somewhat compact textual representation of the
// specified goroutines, by ignoring the often quite lengthy backtrace
// information.
func (matcher *HaveLeakedMatcher) listGoroutines(gs []goroutine.Goroutine, indentation uint) string {
	var buff strings.Builder
	indent := strings.Repeat(format.Indent, int(indentation))
	for _, g := range gs {
		buff.WriteString(indent)
		buff.WriteString(g.String())
	}
	return buff.String()
}

// filter returns a list of leaked goroutines by removing all expected
// goroutines from the given list of goroutines, using the specified checkers.
// The calling goroutine is always filtered out automatically. A checker checks
// if a certain goroutine is expected (then it gets filtered out), or not. If
// all checkers do not signal that they expect a certain goroutine then this
// goroutine is considered to be a leak.
func (matcher *HaveLeakedMatcher) filter(
	goroutines []goroutine.Goroutine, filters []types.GomegaMatcher,
) ([]goroutine.Goroutine, error) {
	gs := make([]goroutine.Goroutine, 0, len(goroutines))
	myID := goroutine.Current().ID
nextgoroutine:
	for _, g := range goroutines {
		if g.ID == myID {
			continue
		}
		for _, filter := range filters {
			matches, err := filter.Match(g)
			if err != nil {
				return nil, err
			}
			if matches {
				continue nextgoroutine
			}
		}
		gs = append(gs, g)
	}
	return gs, nil
}
