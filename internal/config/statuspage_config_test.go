package config_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeStatuspageConfig(t *testing.T, content string) string {
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

func TestLoad_StatuspageDefaults(t *testing.T) {
	path := writeStatuspageConfig(t, "[general]\ninterval = 30\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Statuspage.Enabled {
		t.Error("expected statuspage disabled by default")
	}
	if cfg.Statuspage.BaseURL != "https://api.statuspage.io/v1" {
		t.Errorf("unexpected default base_url: %s", cfg.Statuspage.BaseURL)
	}
}

func TestLoad_StatuspageSection(t *testing.T) {
	path := writeStatuspageConfig(t, `
[statuspage]
enabled = true
api_key = "mykey"
page_id = "pageid123"
component_id = "comp456"
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Statuspage.Enabled {
		t.Error("expected statuspage enabled")
	}
	if cfg.Statuspage.APIKey != "mykey" {
		t.Errorf("unexpected api_key: %s", cfg.Statuspage.APIKey)
	}
	if cfg.Statuspage.PageID != "pageid123" {
		t.Errorf("unexpected page_id: %s", cfg.Statuspage.PageID)
	}
	if cfg.Statuspage.ComponentID != "comp456" {
		t.Errorf("unexpected component_id: %s", cfg.Statuspage.ComponentID)
	}
}

func TestLoad_StatuspageMissingAPIKey(t *testing.T) {
	path := writeStatuspageConfig(t, `
[statuspage]
enabled = true
page_id = "pageid123"
component_id = "comp456"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing api_key")
	}
}

func TestLoad_StatuspageMissingComponentID(t *testing.T) {
	path := writeStatuspageConfig(t, `
[statuspage]
enabled = true
api_key = "mykey"
page_id = "pageid123"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing component_id")
	}
}
