package config

import (
	"os"
	"testing"
)

func writeCustomEventConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	_ = f.Close()
	return f.Name()
}

func TestLoad_CustomEventDefaults(t *testing.T) {
	cfg := defaultCustomEventConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Method != "POST" {
		t.Errorf("expected Method=POST, got %q", cfg.Method)
	}
	if cfg.TimeoutMs != 5000 {
		t.Errorf("expected TimeoutMs=5000, got %d", cfg.TimeoutMs)
	}
}

func TestLoad_CustomEventSection(t *testing.T) {
	path := writeCustomEventConfig(t, `
[custom_event]
enabled = true
url = "https://example.com/hook"
method = "PUT"
timeout_ms = 3000
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	ce := cfg.CustomEvent
	if !ce.Enabled {
		t.Error("expected Enabled=true")
	}
	if ce.URL != "https://example.com/hook" {
		t.Errorf("unexpected URL: %q", ce.URL)
	}
	if ce.Method != "PUT" {
		t.Errorf("unexpected Method: %q", ce.Method)
	}
	if ce.TimeoutMs != 3000 {
		t.Errorf("unexpected TimeoutMs: %d", ce.TimeoutMs)
	}
}

func TestLoad_CustomEventMissingURL(t *testing.T) {
	path := writeCustomEventConfig(t, `
[custom_event]
enabled = true
method = "POST"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestLoad_CustomEventInvalidMethod(t *testing.T) {
	path := writeCustomEventConfig(t, `
[custom_event]
enabled = true
url = "https://example.com/hook"
method = "DELETE"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for unsupported method")
	}
}
