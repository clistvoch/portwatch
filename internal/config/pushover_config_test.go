package config

import (
	"os"
	"testing"
)

func writePushoverConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "pushover_*.toml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	_ = f.Close()
	return f.Name()
}

func TestLoad_PushoverDefaults(t *testing.T) {
	cfg := defaultPushoverConfig()
	if cfg.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.Priority != 0 {
		t.Errorf("expected default priority 0, got %d", cfg.Priority)
	}
	if cfg.Sound != "pushover" {
		t.Errorf("expected default sound 'pushover', got %s", cfg.Sound)
	}
	if cfg.Title != "portwatch alert" {
		t.Errorf("expected default title 'portwatch alert', got %s", cfg.Title)
	}
}

func TestLoad_PushoverSection(t *testing.T) {
	path := writePushoverConfig(t, `
[pushover]
enabled = true
api_token = "tok123"
user_key = "usr456"
title = "my alerts"
priority = 1
sound = "siren"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := cfg.Pushover
	if !p.Enabled {
		t.Error("expected enabled=true")
	}
	if p.APIToken != "tok123" {
		t.Errorf("expected api_token 'tok123', got %s", p.APIToken)
	}
	if p.UserKey != "usr456" {
		t.Errorf("expected user_key 'usr456', got %s", p.UserKey)
	}
	if p.Priority != 1 {
		t.Errorf("expected priority 1, got %d", p.Priority)
	}
	if p.Sound != "siren" {
		t.Errorf("expected sound 'siren', got %s", p.Sound)
	}
}

func TestLoad_PushoverMissingAPIKey(t *testing.T) {
	path := writePushoverConfig(t, `
[pushover]
enabled = true
user_key = "usr456"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing api_token")
	}
}

func TestLoad_PushoverInvalidPriority(t *testing.T) {
	path := writePushoverConfig(t, `
[pushover]
enabled = true
api_token = "tok123"
user_key = "usr456"
priority = 5
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid priority")
	}
}
