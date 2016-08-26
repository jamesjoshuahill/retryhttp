package retryhttp

import (
	"math"
	"time"
)

type AnotherRetryPolicy struct {
	MaxBackOff time.Duration
	MaxDelay   time.Duration
}

func (policy AnotherRetryPolicy) DelayFor(attempt uint) (time.Duration, bool) {
	if policy.sumDelaysIncluding(attempt) > policy.MaxBackOff {
		return 0, false
	}
	return policy.delayFor(attempt), true
}

func (policy AnotherRetryPolicy) delayFor(attempt uint) time.Duration {
	seconds := math.Pow(2, float64(attempt-1))
	delay := time.Duration(seconds) * time.Second
	if delay > policy.MaxDelay {
		return policy.MaxDelay
	}
	return delay
}

func (policy AnotherRetryPolicy) sumDelaysIncluding(currentAttempt uint) (delay time.Duration) {
	for attempt := currentAttempt; attempt > 0; attempt-- {
		delay += policy.delayFor(attempt)
	}
	return
}
