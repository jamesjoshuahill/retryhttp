package retryhttp

import (
	"math"
	"time"
)

//go:generate counterfeiter . RetryPolicy

type RetryPolicy interface {
	DelayFor(uint) (time.Duration, bool)
}

type ExponentialRetryPolicy struct {
	MaxBackOff              time.Duration
	MaxDelayBetweenAttempts time.Duration
}

func (policy ExponentialRetryPolicy) DelayFor(attempt uint) (time.Duration, bool) {
	if policy.sumDelaysIncluding(attempt) > policy.MaxBackOff {
		return 0, false
	}
	return policy.delayFor(attempt), true
}

func (policy ExponentialRetryPolicy) delayFor(attempt uint) time.Duration {
	seconds := math.Pow(2, float64(attempt-1))
	delay := time.Duration(seconds) * time.Second
	if delay > policy.MaxDelayBetweenAttempts {
		return policy.MaxDelayBetweenAttempts
	}
	return delay
}

func (policy ExponentialRetryPolicy) sumDelaysIncluding(currentAttempt uint) (delay time.Duration) {
	for attempt := currentAttempt; attempt > 0; attempt-- {
		delay += policy.delayFor(attempt)
	}
	return
}
