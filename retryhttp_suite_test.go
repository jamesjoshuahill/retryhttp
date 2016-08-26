package retryhttp_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRetryhttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Retryhttp Suite")
}

type row struct {
	attempts     uint
	delay        time.Duration
	keepRetrying bool
}
