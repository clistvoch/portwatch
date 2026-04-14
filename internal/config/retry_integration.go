package config

import (
	"time"

	"github.com/patrickdappollonio/portwatch/internal/alert"
)

// parseDuration is a thin wrapper used by retry and other duration-based configs.
func parseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

// RetryConfigFromSettings converts a WebhookRetryConfig into the alert.RetryConfig
// used at runtime. Durations are parsed; invalid values fall back to defaults.
func RetryConfigFromSettings(c WebhookRetryConfig) alert.RetryConfig {
	def := alert.DefaultRetryConfig()
	initial, err := parseDuration(c.InitialDelay)
	if err != nil {
		initial = def.InitialDelay
	}
	maxD, err := parseDuration(c.MaxDelay)
	if err != nil {
		maxD = def.MaxDelay
	}
	attempts := c.MaxAttempts
	if attempts < 1 {
		attempts = def.MaxAttempts
	}
	mult := c.Multiplier
	if mult < 1.0 {
		mult = def.Multiplier
	}
	return alert.RetryConfig{
		MaxAttempts:  attempts,
		InitialDelay: initial,
		MaxDelay:     maxD,
		Multiplier:   mult,
	}
}
