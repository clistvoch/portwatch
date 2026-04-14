package config

import (
	"testing"
)

func TestWebhookRetryConfig_Defaults(t *testing.T) {
	c := defaultWebhookRetryConfig()
	if c.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", c.MaxAttempts)
	}
	if c.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", c.Multiplier)
	}
	if c.InitialDelay != "500ms" {
		t.Errorf("expected InitialDelay=500ms, got %s", c.InitialDelay)
	}
	if c.MaxDelay != "10s" {
		t.Errorf("expected MaxDelay=10s, got %s", c.MaxDelay)
	}
}

func TestValidateWebhookRetry_Valid(t *testing.T) {
	c := defaultWebhookRetryConfig()
	if err := validateWebhookRetry(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateWebhookRetry_InvalidAttempts(t *testing.T) {
	c := defaultWebhookRetryConfig()
	c.MaxAttempts = 0
	if err := validateWebhookRetry(c); err == nil {
		t.Fatal("expected error for MaxAttempts=0")
	}
}

func TestValidateWebhookRetry_InvalidMultiplier(t *testing.T) {
	c := defaultWebhookRetryConfig()
	c.Multiplier = 0.5
	if err := validateWebhookRetry(c); err == nil {
		t.Fatal("expected error for Multiplier<1")
	}
}

func TestValidateWebhookRetry_BadInitialDelay(t *testing.T) {
	c := defaultWebhookRetryConfig()
	c.InitialDelay = "notaduration"
	if err := validateWebhookRetry(c); err == nil {
		t.Fatal("expected error for bad initial_delay")
	}
}
