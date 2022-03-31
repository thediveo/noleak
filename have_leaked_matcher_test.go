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
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/thediveo/noleak/goroutine"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HaveLeaked", func() {

	It("considers testing and runtime goroutines not to be leaks", func() {
		Expect(Goroutines()).NotTo(HaveLeaked(), "should not find any leaks by default")
	})

	When("using signals", func() {

		It("doesn't find leaks", func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			Eventually(Goroutines).WithTimeout(2*time.Second).WithPolling(250*time.Millisecond).
				ShouldNot(HaveLeaked(), "found signal.Notify leaks")

			signal.Reset(os.Interrupt)
			Eventually(Goroutines).WithTimeout(2*time.Second).WithPolling(250*time.Millisecond).
				ShouldNot(HaveLeaked(), "found signal.Reset leaks")
		})

	})

	It("checks against list of expected goroutines", func() {
		By("taking a snapshot")
		gs := Goroutines()
		m := HaveLeaked(gs)

		By("starting a goroutine")
		done := make(chan struct{})
		var once sync.Once
		go func() {
			<-done
		}()
		defer once.Do(func() { close(done) })

		By("detecting the goroutine")
		Expect(m.Match(Goroutines())).To(BeTrue())

		By("terminating the goroutine and ensuring it has terminated")
		once.Do(func() { close(done) })
		Eventually(func() (bool, error) {
			return m.Match(Goroutines())
		}).Should(BeFalse())
	})

	Context("failure messages", func() {

		var snapshot []goroutine.Goroutine
		var done chan struct{}

		BeforeEach(func() {
			snapshot = Goroutines()
			done = make(chan struct{})
			go func() {
				<-done
			}()
		})

		AfterEach(func() {
			close(done)
			Eventually(Goroutines).ShouldNot(HaveLeaked(snapshot))
		})

		It("returns a failure message", func() {
			m := HaveLeaked(snapshot)
			gs := Goroutines()
			Expect(m.Match(gs)).To(BeTrue())
			Expect(m.FailureMessage(gs)).To(MatchRegexp(
				`Expected to leak goroutines:\n    Goroutine ID: \d+, state: .*, top function: .*`))
		})

		It("returns a negated failure message", func() {
			m := HaveLeaked(snapshot)
			gs := Goroutines()
			Expect(m.Match(gs)).To(BeTrue())
			Expect(m.NegatedFailureMessage(gs)).To(MatchRegexp(
				`Expected not to leak goroutines:\n    Goroutine ID: \d+, state: .*, top function: .*`))

		})

		When("things go wrong", func() {

			It("rejects unsupported filter args types", func() {
				Expect(func() { _ = HaveLeaked(42) }).To(PanicWith(
					"HaveLeaked expected a string, []Goroutine, or GomegaMatcher, but got:\n    <int>: 42"))
			})

			It("accepts plain strings as filters", func() {
				m := HaveLeaked("foo.bar")
				Expect(m.Match([]goroutine.Goroutine{
					{TopFunction: "foo.bar"},
				})).To(BeFalse())
			})

			It("expects actual to be a slice of goroutine.Goroutine", func() {
				m := HaveLeaked()
				Expect(m.Match(nil)).Error().To(MatchError(
					"HaveLeaked matcher expects an array or slice of goroutines.  Got:\n    <nil>: nil"))
				Expect(m.Match("foo!")).Error().To(MatchError(
					"HaveLeaked matcher expects an array or slice of goroutines.  Got:\n    <string>: foo!"))
				Expect(m.Match([]string{"foo!"})).Error().To(MatchError(
					"HaveLeaked matcher expects an array or slice of goroutines.  Got:\n    <[]string | len:1, cap:1>: [\"foo!\"]"))
			})

			It("handles filter matcher errors", func() {
				m := HaveLeaked(HaveField("foobar", BeNil()))
				Expect(m.Match([]goroutine.Goroutine{
					{ID: 0},
				})).Error().To(HaveOccurred())
			})

		})

	})

	Context("wrapped around test nodes", func() {

		var snapshot []goroutine.Goroutine

		When("not leaking", func() {

			BeforeEach(func() {
				snapshot = Goroutines()
			})

			AfterEach(func() {
				Eventually(Goroutines).ShouldNot(HaveLeaked(snapshot))
			})

			It("doesn't leak in test", func() {
				// nothing
			})

		})

		When("leaking", func() {

			done := make(chan struct{})

			BeforeEach(func() {
				snapshot = Goroutines()
			})

			AfterEach(func() {
				Expect(Goroutines()).To(HaveLeaked(snapshot))
				close(done)
				Eventually(Goroutines).ShouldNot(HaveLeaked(snapshot))
			})

			It("leaks in test", func() {
				go func() {
					<-done
				}()
			})

		})

	})

})
