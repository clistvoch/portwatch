package config

import (
	"os"
	"testing"
)

func writeSplunkConfig(t *testing.T, content string) string {
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

func TestLoad_SplunkDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Splunk.Enabled {
		t.Error("expected splunk disabled by default")
	}
	if cfg.Splunk.Index != "main" {
		t.Errorf("expected index 'main', got %q", cfg.Splunk.Index)
	}
	if cfg.Splunk.SourceType != "portwatch" {
		t.Errorf("expected source_type 'portwatch', got %q", cfg.Splunk.SourceType)
	}
	if cfg.Splunk.Timeout != 10 {
		t.Errorf("expected timeout 10, got %d", cfg.Splunk.Timeout)
	}
}

func TestLoad_SplunkSection(t *testing.T) {
	path := writeSplunkConfig(t, `
[splunk]
enabled = true
url = "http://splunk.example.com:8088/services/collector"
token = "abc-123"
index = "portwatch"
source_type = "json"
timeout_seconds = 5
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Splunk.Enabled {
		t.Error("expected splunk enabled")
	}
	if cfg.Splunk.URL != "http://splunk.example.com:8088/services/collector" {
		t.Errorf("unexpected url: %s", cfg.Splunk.URL)
	}
	if cfg.Splunk.Token != "abc-123" {
		t.Errorf("unexpected token: %s", cfg.Splunk.Token)
	}
	if cfg.Splunk.Index != "portwatch" {
		t.Errorf("unexpected index: %s", cfg.Splunk.Index)
	}
	if cfg.Splunk.Timeout != 5 {
		t.Errorf("unexpected timeout: %d", cfg.Splunk.Timeout)
	}
}

func TestLoad_SplunkMissingURL(t *testing.T) {
	path := writeSplunkConfig(t, `
[splunk]
enabled = true
token = "abc-123"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing url")
	}
}

func TestLoad_SplunkMissingToken(t *testing.T) {
	path := writeSplunkConfig(t, `
[splunk]
enabled = true
url = "http://splunk.example.com:8088/services/collector"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing token")
	}
}
