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
	MaxBackOff time.Duration
}

const maxRetryDelay = 16 * time.Second

func (policy ExponentialRetryPolicy) DelayFor(attempt uint) (time.Duration, bool) {
	if sumDelaysIncluding(attempt) > policy.MaxBackOff {
		return 0, false
	}
	return delayFor(attempt), true
}

func delayFor(attempt uint) time.Duration {
	seconds := math.Pow(2, float64(attempt-1))
	delay := time.Duration(seconds) * time.Second
	if delay > maxRetryDelay {
		return maxRetryDelay
	}
	return delay
}

func sumDelaysIncluding(currentAttempt uint) (delay time.Duration) {
	for attempt := currentAttempt; attempt > 0; attempt-- {
		delay += delayFor(attempt)
	}
	return
}
