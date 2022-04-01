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

var _ = Describe("IgnoringGoroutines matcher", func() {

	It("returns an error for an invalid actual", func() {
		m := IgnoringGoroutines(Goroutines())
		Expect(m.Match(nil)).Error().To(MatchError(
			"IgnoringGoroutines matcher expects a goroutine.Goroutine or *goroutine.Goroutine.  Got:\n    <nil>: nil"))
	})

	It("matches", func() {
		gs := Goroutines()
		me := gs[0]
		m := IgnoringGoroutines(gs)
		Expect(m.Match(me)).To(BeTrue())
		Expect(m.Match(gs[1])).To(BeTrue())
		Expect(m.Match(goroutine.Goroutine{})).To(BeFalse())
	})

	It("returns failure messages", func() {
		m := IgnoringGoroutines(Goroutines())
		Expect(m.FailureMessage(goroutine.Goroutine{})).To(MatchRegexp(
			`Expected\n    <goroutine.Goroutine>: {ID: 0, State: "", TopFunction: "", CreatorFunction: "", CreatorLocation: ""}\nto be contained in the list of expected goroutine IDs\n    <\[\]uint64 | len:\d+, cap:\d+>: [.*]`))
		Expect(m.NegatedFailureMessage(goroutine.Goroutine{})).To(MatchRegexp(
			`Expected\n    <goroutine.Goroutine>: {ID: 0, State: "", TopFunction: "", CreatorFunction: "", CreatorLocation: ""}\nnot to be contained in the list of expected goroutine IDs\n    <\[\]uint64 | len:\d+, cap:\d+>: [.*]`))
	})

})
