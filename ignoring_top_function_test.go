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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/thediveo/noleak/goroutine"
)

var _ = Describe("IgnoringTopFunction matcher", func() {

	It("returns an error for an invalid actual", func() {
		m := IgnoringTopFunction("foo.bar")
		Expect(m.Match(nil)).Error().To(MatchError("IgnoringTopFunction matcher expects a goroutine.Goroutine or *goroutine.Goroutine.  Got:\n    <nil>: nil"))
	})

	It("matches a toplevel function by full name", func() {
		m := IgnoringTopFunction("foo.bar")
		Expect(m.Match(goroutine.Goroutine{
			TopFunction: "foo.bar",
		})).To(BeTrue())
		Expect(m.Match(goroutine.Goroutine{
			TopFunction: "main.main",
		})).To(BeFalse())
	})

	It("matches a toplevel function by prefix", func() {
		m := IgnoringTopFunction("foo...")
		Expect(m.Match(goroutine.Goroutine{
			TopFunction: "foo.bar",
		})).To(BeTrue())
		Expect(m.Match(goroutine.Goroutine{
			TopFunction: "foo",
		})).To(BeFalse())
		Expect(m.Match(goroutine.Goroutine{
			TopFunction: "spanish.inquisition",
		})).To(BeFalse())
	})

	It("matches a toplevel function by prefix", func() {
		m := IgnoringTopFunction("foo.bar [worried]")
		Expect(m.Match(goroutine.Goroutine{
			TopFunction: "foo.bar",
			State:       "worried, stalled",
		})).To(BeTrue())
		Expect(m.Match(goroutine.Goroutine{
			TopFunction: "foo.bar",
			State:       "uneasy, anxious",
		})).To(BeFalse())
	})

	It("returns failure messages", func() {
		m := IgnoringTopFunction("foo.bar")
		Expect(m.FailureMessage(goroutine.Goroutine{ID: 42, TopFunction: "foo"})).To(Equal(
			"Expected\n    <goroutine.Goroutine>: {ID: 42, State: \"\", TopFunction: \"foo\", Backtrace: \"\"}\nto have the topmost function \"foo.bar\""))
		Expect(m.NegatedFailureMessage(goroutine.Goroutine{ID: 42, TopFunction: "foo"})).To(Equal(
			"Expected\n    <goroutine.Goroutine>: {ID: 42, State: \"\", TopFunction: \"foo\", Backtrace: \"\"}\nnot\n    <string>: to have the topmost function \"foo.bar\""))

		m = IgnoringTopFunction("foo.bar [worried]")
		Expect(m.FailureMessage(goroutine.Goroutine{ID: 42, TopFunction: "foo"})).To(Equal(
			"Expected\n    <goroutine.Goroutine>: {ID: 42, State: \"\", TopFunction: \"foo\", Backtrace: \"\"}\nto have the topmost function \"foo.bar\" and the state \"worried\""))

		m = IgnoringTopFunction("foo...")
		Expect(m.FailureMessage(goroutine.Goroutine{ID: 42, TopFunction: "foo"})).To(Equal(
			"Expected\n    <goroutine.Goroutine>: {ID: 42, State: \"\", TopFunction: \"foo\", Backtrace: \"\"}\nto have the prefix \"foo.\" for its topmost function"))
	})

})
