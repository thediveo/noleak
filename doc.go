/*

Package noleak complements the Gingko/Gomega testing and matchers framework with
matchers for Goroutine leakage detection.

Basics of noleak

To start with,

  Goroutines()

returns information about all (non-dead) goroutines at a particular moment. This
is useful to capture a known correct snapshot and then later taking a new
snapshot and comparing these two snapshots for leaked goroutines.

Next, the matcher

  HaveLeaked()

filters out well-known and expected goroutines from an actual list of goroutines
(passed from Eventually or Expect), hopefully ending up with an empty list of
leaked goroutines. If there are still goroutines left after filtering out the
well-known and expected goroutines, then HaveLeaked() will succeed. Which
actually is usually considered to be failure, so rather to be "suckcess" because
no one wants leaked goroutines.

A typical pattern to detect goroutines leaked in tests is as follows:

   var snapshot []goroutine.Goroutine

   BeforeEach(func() {
	   snapshot = Goroutines()
   })

   AfterEach(func() {
	   // Note: it's "Goroutines", but not "Goroutines()", when using with Eventually!
	   Eventually(Goroutines).ShouldNot(HaveLeaked(snapshot))
   })

Acknowledgement

noleak has been heavily inspired by the Goroutine leak detector
github.com/uber-go/goleak. It's definitely a fine piece of work!

But in the end, we had to decide against trying to hammering down uber-go/goleak
into the Gomega TDD matcher ecosystem, because reusing and wrapping would have
become very awkward. The main reason is that goleak.Find combines all the
different elements of getting actual goroutines information, filtering them,
arriving at a leak conclusion, and even retrying multiple times in one single
exported function. Unfortunately, goleak makes gathering information about all
goroutines an internal matter, so we cannot reuse such functionality elsewhere
outside goleak.Find.

Users of the Gomega ecosystem are already experienced in arriving at conclusions
and retrying temporarily failing expectations: Gomega does it in form of
Eventually().ShouldNot(), and (without the trying aspect) with Expect().NotTo().
So what is missing is only a goroutine leak detector in form of the HaveLeaked
matcher, as well as the ability to specify goroutine filters in order to sort
out the non-leaking (and therefore expected) goroutines, using a few filter
criteria. That is, a few new goroutine-related matchers. In this architecture,
even existing Gomega matchers can optionally be (re)used as the need arises.

References

https://github.com/onsi/gomega and https://github.com/onsi/ginkgo.

*/
package noleak
