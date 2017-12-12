package typederrs

import (
	"time"
	"github.com/jpillora/backoff"
)

type RetryConfig struct {
	backoff *backoff.Backoff
	checker IsRetryableErrChecker
}

type RetryOption func(*RetryConfig)

func RetryWithMinBackoff(d time.Duration) RetryOption {
	return func(b *RetryConfig) {
		b.backoff.Min = d
	}
}

func RetryWithMaxBackoff(d time.Duration) RetryOption {
	return func(b *RetryConfig) {
		b.backoff.Max = d
	}
}

func RetryWithBackoffFactor(f float64) RetryOption {
	return func(b *RetryConfig) {
		b.backoff.Factor = f
	}
}

func RetryWithBackoffJitter(t bool) RetryOption {
	return func(b *RetryConfig) {
		b.backoff.Jitter = t
	}
}

func RetryWithRetryableErrChecker(ch IsRetryableErrChecker) RetryOption {
	return func(b *RetryConfig) {
		b.checker = ch
	}
}

func DoWithRetries(doer func() error, opts ...RetryOption) error {

	conf := RetryConfig{
		backoff: &backoff.Backoff{Min: 2 * time.Second, Max: 5 * time.Minute},
		checker: &RetryableErrCheck{},
	}
	for _, f := range opts {
		f(&conf)
	}

	var err error

	for numRetries := 0; numRetries < 5; numRetries++ {
		err = doer()
		if err == nil {
			return nil
		}
		if !conf.checker.IsRetryableError(err) {
			return err
		}
		time.Sleep(conf.backoff.Duration())
	}

	return Newf("too many retries: %v", err)
}
