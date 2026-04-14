package config

import "fmt"

// WebhookRetryConfig holds retry/backoff settings for outbound HTTP handlers.
type WebhookRetryConfig struct {
	MaxAttempts int     `toml:"max_attempts"`
	InitialDelay string  `toml:"initial_delay"`
	MaxDelay     string  `toml:"max_delay"`
	Multiplier   float64 `toml:"multiplier"`
}

func defaultWebhookRetryConfig() WebhookRetryConfig {
	return WebhookRetryConfig{
		MaxAttempts:  3,
		InitialDelay: "500ms",
		MaxDelay:     "10s",
		Multiplier:   2.0,
	}
}

func validateWebhookRetry(c WebhookRetryConfig) error {
	if c.MaxAttempts < 1 {
		return ValidationError{Field: "retry.max_attempts", Msg: "must be at least 1"}
	}
	if c.Multiplier < 1.0 {
		return ValidationError{Field: "retry.multiplier", Msg: "must be >= 1.0"}
	}
	if _, err := parseDuration(c.InitialDelay); err != nil {
		return ValidationError{Field: "retry.initial_delay", Msg: fmt.Sprintf("invalid duration: %v", err)}
	}
	if _, err := parseDuration(c.MaxDelay); err != nil {
		return ValidationError{Field: "retry.max_delay", Msg: fmt.Sprintf("invalid duration: %v", err)}
	}
	return nil
}
