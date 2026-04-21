package config

import (
	"testing"
)

func TestGoogleChatHandlerConfig_Defaults(t *testing.T) {
	cfg := defaultGoogleChatHandlerConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.WebhookURL != "" {
		t.Errorf("expected empty WebhookURL, got %q", cfg.WebhookURL)
	}
}

func TestValidateGoogleChatHandler_Valid(t *testing.T) {
	cfg := GoogleChatHandlerConfig{
		Enabled:    true,
		WebhookURL: "https://chat.googleapis.com/v1/spaces/ABC/messages?key=xyz",
	}
	if err := validateGoogleChatHandler(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateGoogleChatHandler_DisabledSkipsValidation(t *testing.T) {
	cfg := GoogleChatHandlerConfig{Enabled: false}
	if err := validateGoogleChatHandler(cfg); err != nil {
		t.Fatalf("expected no error when disabled, got: %v", err)
	}
}

func TestValidateGoogleChatHandler_MissingWebhook(t *testing.T) {
	cfg := GoogleChatHandlerConfig{Enabled: true, WebhookURL: ""}
	err := validateGoogleChatHandler(cfg)
	if err == nil {
		t.Fatal("expected error for missing webhook_url")
	}
}

func TestGoogleChatHandlerConfigFromSettings_Valid(t *testing.T) {
	s := map[string]string{
		"webhook_url": "https://chat.googleapis.com/v1/spaces/XYZ/messages?key=abc",
	}
	cfg := GoogleChatHandlerConfigFromSettings(s)
	if !cfg.Enabled {
		t.Error("expected Enabled to be true when webhook_url is set")
	}
	if cfg.WebhookURL != s["webhook_url"] {
		t.Errorf("expected WebhookURL %q, got %q", s["webhook_url"], cfg.WebhookURL)
	}
}

func TestGoogleChatHandlerConfigFromSettings_ExplicitDisable(t *testing.T) {
	s := map[string]string{
		"webhook_url": "https://chat.googleapis.com/v1/spaces/XYZ/messages?key=abc",
		"enabled":     "false",
	}
	cfg := GoogleChatHandlerConfigFromSettings(s)
	if cfg.Enabled {
		t.Error("expected Enabled to be false when explicitly set to false")
	}
}
