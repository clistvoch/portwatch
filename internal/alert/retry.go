package alert

import (
	"context"
	"math"
	"time"
)

// RetryConfig controls exponential-backoff retry behaviour.
type RetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultRetryConfig returns sensible retry defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// Do calls fn up to cfg.MaxAttempts times with exponential backoff.
// It returns the last error if all attempts fail, or nil on success.
func Do(ctx context.Context, cfg RetryConfig, fn func() error) error {
	delay := cfg.InitialDelay
	var err error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err = fn(); err == nil {
			return nil
		}
		if attempt == cfg.MaxAttempts-1 {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		next := time.Duration(math.Round(float64(delay) * cfg.Multiplier))
		if next > cfg.MaxDelay {
			next = cfg.MaxDelay
		}
		delay = next
	}
	return err
}
