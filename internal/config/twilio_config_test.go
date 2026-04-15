package config

import (
	"os"
	"testing"
)

func writeTwilioConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_TwilioDefaults(t *testing.T) {
	path := writeTwilioConfig(t, "[scan]\nports = \"80-81\"\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Twilio.Enabled {
		t.Error("expected twilio disabled by default")
	}
	if cfg.Twilio.BaseURL != "https://api.twilio.com" {
		t.Errorf("unexpected base_url: %s", cfg.Twilio.BaseURL)
	}
}

func TestLoad_TwilioSection(t *testing.T) {
	path := writeTwilioConfig(t, `
[scan]
ports = "80-81"

[twilio]
enabled = true
account_sid = "ACtest123"
auth_token = "secret"
from = "+15550001111"
to = "+15559998888"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Twilio.Enabled {
		t.Error("expected twilio enabled")
	}
	if cfg.Twilio.AccountSID != "ACtest123" {
		t.Errorf("unexpected account_sid: %s", cfg.Twilio.AccountSID)
	}
	if cfg.Twilio.To != "+15559998888" {
		t.Errorf("unexpected to: %s", cfg.Twilio.To)
	}
}

func TestLoad_TwilioMissingAccountSID(t *testing.T) {
	path := writeTwilioConfig(t, `
[scan]
ports = "80-81"

[twilio]
enabled = true
auth_token = "secret"
from = "+15550001111"
to = "+15559998888"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing account_sid")
	}
}

func TestLoad_TwilioMissingTo(t *testing.T) {
	path := writeTwilioConfig(t, `
[scan]
ports = "80-81"

[twilio]
enabled = true
account_sid = "ACtest123"
auth_token = "secret"
from = "+15550001111"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing to number")
	}
}
