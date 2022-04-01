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

var _ = Describe("IgnoringInBacktrace matcher", func() {

	It("returns an error for an invalid actual", func() {
		m := IgnoringInBacktrace("foo.bar")
		Expect(m.Match(nil)).Error().To(MatchError(
			"IgnoringInBacktrace matcher expects a goroutine.Goroutine or *goroutine.Goroutine.  Got:\n    <nil>: nil"))
	})

	It("matches", func() {
		m := IgnoringInBacktrace("github.com/thediveo/noleak/goroutine.stacks")
		Expect(m.Match(somefunction())).To(BeTrue())
	})

	It("returns failure messages", func() {
		m := IgnoringInBacktrace("foo.bar")
		Expect(m.FailureMessage(goroutine.Goroutine{Backtrace: "abc"})).To(MatchRegexp(
			`Expected\n    <goroutine.Goroutine>: {ID: 0, State: "", TopFunction: "", CreatorFunction: "", CreatorLocation: ""}\nto contain "foo.bar" in the goroutine's stack backtrace`))
		Expect(m.NegatedFailureMessage(goroutine.Goroutine{Backtrace: "abc"})).To(MatchRegexp(
			`Expected\n    <goroutine.Goroutine>: {ID: 0, State: "", TopFunction: "", CreatorFunction: "", CreatorLocation: ""}\nnot to contain "foo.bar" in the goroutine's stack backtrace`))
	})

})

func somefunction() goroutine.Goroutine {
	return goroutine.Current()
}
