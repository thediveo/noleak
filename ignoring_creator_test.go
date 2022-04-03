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
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/thediveo/noleak/goroutine"
)

func creator() goroutine.Goroutine {
	ch := make(chan goroutine.Goroutine)
	go func() {
		ch <- goroutine.Current()
	}()
	return <-ch
}

var _ = Describe("IgnoringCreator matcher", func() {

	It("returns an error for an invalid actual", func() {
		m := IgnoringCreator("foo.bar")
		Expect(m.Match(nil)).Error().To(MatchError("IgnoringCreator matcher expects a goroutine.Goroutine or *goroutine.Goroutine.  Got:\n    <nil>: nil"))
	})

	It("matches a creator function by full name", func() {
		type T struct{}
		pkg := reflect.TypeOf(T{}).PkgPath()
		m := IgnoringCreator(pkg + ".creator")
		g := creator()
		Expect(m.Match(g)).To(BeTrue(), "creator %v", g.String())
		Expect(m.Match(goroutine.Current())).To(BeFalse())
	})

	It("matches a toplevel function by prefix", func() {
		type T struct{}
		pkg := reflect.TypeOf(T{}).PkgPath()
		m := IgnoringCreator(pkg + "...")
		g := creator()
		Expect(m.Match(g)).To(BeTrue(), "creator %v", g.String())
		Expect(m.Match(goroutine.Current())).To(BeFalse())
		Expect(m.Match(goroutine.Goroutine{
			TopFunction: "spanish.inquisition",
		})).To(BeFalse())
	})

	It("returns failure messages", func() {
		m := IgnoringCreator("foo.bar")
		Expect(m.FailureMessage(goroutine.Goroutine{ID: 42, TopFunction: "foo"})).To(Equal(
			"Expected\n    <goroutine.Goroutine>: {ID: 42, State: \"\", TopFunction: \"foo\", CreatorFunction: \"\", BornAt: \"\"}\nto be created by \"foo.bar\""))
		Expect(m.NegatedFailureMessage(goroutine.Goroutine{ID: 42, TopFunction: "foo"})).To(Equal(
			"Expected\n    <goroutine.Goroutine>: {ID: 42, State: \"\", TopFunction: \"foo\", CreatorFunction: \"\", BornAt: \"\"}\nnot to be created by \"foo.bar\""))

		m = IgnoringCreator("foo...")
		Expect(m.FailureMessage(goroutine.Goroutine{ID: 42, TopFunction: "foo"})).To(Equal(
			"Expected\n    <goroutine.Goroutine>: {ID: 42, State: \"\", TopFunction: \"foo\", CreatorFunction: \"\", BornAt: \"\"}\nto be created by a function with prefix \"foo.\""))
	})

})
