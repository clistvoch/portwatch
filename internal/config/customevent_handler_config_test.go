package config

import (
	"testing"
)

func TestCustomEventHandlerConfig_Defaults(t *testing.T) {
	cfg := defaultCustomEventHandlerConfig()
	if cfg.Method != "POST" {
		t.Errorf("expected POST, got %s", cfg.Method)
	}
	if cfg.Timeout != 10 {
		t.Errorf("expected timeout 10, got %d", cfg.Timeout)
	}
	if cfg.URL != "" {
		t.Errorf("expected empty URL, got %s", cfg.URL)
	}
}

func TestValidateCustomEventHandler_Valid(t *testing.T) {
	cfg := CustomEventHandlerConfig{
		URL:     "https://example.com/hook",
		Method:  "POST",
		Timeout: 5,
		Headers: map[string]string{},
	}
	if err := validateCustomEventHandler(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCustomEventHandler_MissingURL(t *testing.T) {
	cfg := CustomEventHandlerConfig{Method: "POST", Timeout: 5, Headers: map[string]string{}}
	if err := validateCustomEventHandler(cfg); err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestValidateCustomEventHandler_InvalidMethod(t *testing.T) {
	cfg := CustomEventHandlerConfig{
		URL:     "https://example.com/hook",
		Method:  "DELETE",
		Timeout: 5,
		Headers: map[string]string{},
	}
	if err := validateCustomEventHandler(cfg); err == nil {
		t.Fatal("expected error for unsupported method")
	}
}

func TestCustomEventHandlerConfigFromSettings_Valid(t *testing.T) {
	s := map[string]string{
		"url":            "https://hooks.example.com/notify",
		"method":         "put",
		"timeout_seconds": "15",
		"header_X-Token": "abc123",
	}
	cfg, err := CustomEventHandlerConfigFromSettings(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Method != "PUT" {
		t.Errorf("expected PUT, got %s", cfg.Method)
	}
	if cfg.Timeout != 15 {
		t.Errorf("expected 15, got %d", cfg.Timeout)
	}
	if cfg.Headers["X-Token"] != "abc123" {
		t.Errorf("expected header X-Token=abc123, got %s", cfg.Headers["X-Token"])
	}
}

func TestCustomEventHandlerConfigFromSettings_MissingURL(t *testing.T) {
	s := map[string]string{"method": "POST"}
	_, err := CustomEventHandlerConfigFromSettings(s)
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}
