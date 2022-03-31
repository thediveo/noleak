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

var _ = Describe("utilities", func() {

	Context("G(oroutine) descriptions", func() {

		It("returns an error for actual <nil>", func() {
			Expect(func() { _, _ = G(nil, "foo") }).NotTo(Panic())
			Expect(G(nil, "foo")).Error().To(MatchError("foo matcher expects a goroutine.Goroutine or *goroutine.Goroutine.  Got:\n    <nil>: nil"))
		})

		It("returns an error when passing something that's not a goroutine by any means", func() {
			Expect(func() { _, _ = G("foobar", "foo") }).NotTo(Panic())
			Expect(G("foobar", "foo")).Error().To(MatchError("foo matcher expects a goroutine.Goroutine or *goroutine.Goroutine.  Got:\n    <string>: foobar"))
		})

		It("returns a goroutine", func() {
			actual := goroutine.Goroutine{ID: 42}
			g, err := G(actual, "foo")
			Expect(err).NotTo(HaveOccurred())
			Expect(g.ID).To(Equal(42))

			g, err = G(&actual, "foo")
			Expect(err).NotTo(HaveOccurred())
			Expect(g.ID).To(Equal(42))
		})

	})

	It("returns a list of Goroutine IDs in textual format", func() {
		Expect(goids(nil)).To(BeEmpty())
		Expect(goids([]goroutine.Goroutine{
			{ID: 666},
			{ID: 42},
		})).To(Equal("42, 666"))
		Expect(goids([]goroutine.Goroutine{
			{ID: 42},
		})).To(Equal("42"))
	})

})
