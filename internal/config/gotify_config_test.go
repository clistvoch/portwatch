package config

import (
	"os"
	"testing"
)

func writeGotifyConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_GotifyDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Gotify.Enabled {
		t.Error("expected gotify disabled by default")
	}
	if cfg.Gotify.Priority != 5 {
		t.Errorf("expected default priority 5, got %d", cfg.Gotify.Priority)
	}
}

func TestLoad_GotifySection(t *testing.T) {
	path := writeGotifyConfig(t, `
[gotify]
enabled = true
url = "http://gotify.example.com"
token = "abc123"
priority = 8
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Gotify.Enabled {
		t.Error("expected gotify enabled")
	}
	if cfg.Gotify.URL != "http://gotify.example.com" {
		t.Errorf("unexpected url: %s", cfg.Gotify.URL)
	}
	if cfg.Gotify.Token != "abc123" {
		t.Errorf("unexpected token: %s", cfg.Gotify.Token)
	}
	if cfg.Gotify.Priority != 8 {
		t.Errorf("expected priority 8, got %d", cfg.Gotify.Priority)
	}
}

func TestLoad_GotifyMissingToken(t *testing.T) {
	path := writeGotifyConfig(t, `
[gotify]
enabled = true
url = "http://gotify.example.com"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing token")
	}
}

func TestLoad_GotifyInvalidPriority(t *testing.T) {
	path := writeGotifyConfig(t, `
[gotify]
enabled = true
url = "http://gotify.example.com"
token = "tok"
priority = 99
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for invalid priority")
	}
}
