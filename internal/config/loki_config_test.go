package config

import (
	"os"
	"testing"
)

func writeLokiConfig(t *testing.T, content string) string {
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

func TestLoad_LokiDefaults(t *testing.T) {
	path := writeLokiConfig(t, "[scanner]\nport_range = \"1-1024\"\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Loki.Enabled {
		t.Error("expected loki disabled by default")
	}
	if cfg.Loki.URL != "" {
		t.Errorf("expected empty URL, got %q", cfg.Loki.URL)
	}
	if cfg.Loki.Labels["job"] != "portwatch" {
		t.Errorf("expected default job label, got %q", cfg.Loki.Labels["job"])
	}
}

func TestLoad_LokiSection(t *testing.T) {
	path := writeLokiConfig(t, `
[scanner]
port_range = "1-1024"

[loki]
enabled = true
url = "http://loki.example.com:3100"

[loki.labels]
job = "portwatch"
env = "prod"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Loki.Enabled {
		t.Error("expected loki enabled")
	}
	if cfg.Loki.URL != "http://loki.example.com:3100" {
		t.Errorf("unexpected URL: %s", cfg.Loki.URL)
	}
	if cfg.Loki.Labels["env"] != "prod" {
		t.Errorf("expected env=prod label, got %q", cfg.Loki.Labels["env"])
	}
}

func TestLoad_LokiMissingURL(t *testing.T) {
	path := writeLokiConfig(t, `
[scanner]
port_range = "1-1024"

[loki]
enabled = true
url = ""
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing loki url, got nil")
	}
}
