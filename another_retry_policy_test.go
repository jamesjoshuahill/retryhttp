package retryhttp_test

import (
	"fmt"
	"time"

	"github.com/concourse/retryhttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AnotherRetryPolicy", func() {
	var (
		attempts uint

		policy retryhttp.AnotherRetryPolicy

		delay        time.Duration
		keepRetrying bool
	)

	var testAttempts = func(rows []row) {
		for _, row := range rows {
			row := row

			Context(fmt.Sprintf("after %d failed attempts", row.attempts), func() {
				BeforeEach(func() {
					attempts = row.attempts
				})

				It(fmt.Sprintf("returns a %s delay", row.delay), func() {
					Expect(delay).To(Equal(row.delay))
				})

				if row.keepRetrying {
					It("keeps retrying", func() {
						Expect(keepRetrying).To(BeTrue())
					})
				} else {
					It("gives up", func() {
						Expect(keepRetrying).To(BeFalse())
					})
				}
			})
		}
	}

	JustBeforeEach(func() {
		delay, keepRetrying = policy.DelayFor(attempts)
	})

	Context("when the max delay between attempts is 16 seconds", func() {
		Context("with a 1 second timeout", func() {
			BeforeEach(func() {
				policy = retryhttp.AnotherRetryPolicy{
					MaxBackOff: 1 * time.Second,
					MaxDelay:   16 * time.Second,
				}
			})

			testAttempts([]row{
				{0, 0 * time.Second, true},
				{1, 1 * time.Second, true},
				{2, 0, false},
				{3, 0, false},
			})
		})

		Context("with a 3 second timeout", func() {
			BeforeEach(func() {
				policy = retryhttp.AnotherRetryPolicy{
					MaxBackOff: 3 * time.Second,
					MaxDelay:   16 * time.Second,
				}
			})

			testAttempts([]row{
				{0, 0 * time.Second, true},
				{1, 1 * time.Second, true},
				{2, 2 * time.Second, true},
				{3, 0, false},
				{4, 0, false},
			})
		})

		Context("with a 7 second timeout", func() {
			BeforeEach(func() {
				policy = retryhttp.AnotherRetryPolicy{
					MaxBackOff: 7 * time.Second,
					MaxDelay:   16 * time.Second,
				}
			})

			testAttempts([]row{
				{0, 0 * time.Second, true},
				{1, 1 * time.Second, true},
				{2, 2 * time.Second, true},
				{3, 4 * time.Second, true},
				{4, 0, false},
				{5, 0, false},
			})
		})

		Context("with a 5 minute timeout", func() {
			BeforeEach(func() {
				policy = retryhttp.AnotherRetryPolicy{
					MaxBackOff: 5 * time.Minute,
					MaxDelay:   16 * time.Second,
				}
			})

			testAttempts([]row{
				{0, 0 * time.Second, true},
				{1, 1 * time.Second, true},
				{2, 2 * time.Second, true},
				{3, 4 * time.Second, true},
				{4, 8 * time.Second, true},
				{5, 16 * time.Second, true},
				{6, 16 * time.Second, true},
				{7, 16 * time.Second, true},
				{8, 16 * time.Second, true},
				{9, 16 * time.Second, true},
				{20, 16 * time.Second, true},
				{21, 16 * time.Second, true},
				{22, 0, false},
				{23, 0, false},
			})
		})
	})

	Context("when the max delay between attempts is 30 seconds", func() {
		Context("with a 5 minute timeout", func() {
			BeforeEach(func() {
				policy = retryhttp.AnotherRetryPolicy{
					MaxBackOff: 5 * time.Minute,
					MaxDelay:   30 * time.Second,
				}
			})

			testAttempts([]row{
				{0, 0 * time.Second, true},
				{1, 1 * time.Second, true},
				{2, 2 * time.Second, true},
				{3, 4 * time.Second, true},
				{4, 8 * time.Second, true},
				{5, 16 * time.Second, true},
				{6, 30 * time.Second, true},
				{7, 30 * time.Second, true},
				{8, 30 * time.Second, true},
				{9, 30 * time.Second, true},
				{20, 0, false},
				{21, 0, false},
			})
		})
	})
})
