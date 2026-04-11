package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/portwatch/internal/config"
)

func writeOpsGenieConfig(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.toml")
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_OpsGenieDefaults(t *testing.T) {
	p := writeOpsGenieConfig(t, "")
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OpsGenie.Enabled {
		t.Error("expected opsgenie to be disabled by default")
	}
	if cfg.OpsGenie.Priority != "P3" {
		t.Errorf("expected default priority P3, got %q", cfg.OpsGenie.Priority)
	}
	if cfg.OpsGenie.APIBaseURL != "https://api.opsgenie.com" {
		t.Errorf("unexpected default api_base_url: %q", cfg.OpsGenie.APIBaseURL)
	}
}

func TestLoad_OpsGenieSection(t *testing.T) {
	body := `
[opsgenie]
enabled = true
api_key = "secret-key"
team = "platform"
priority = "P2"
`
	p := writeOpsGenieConfig(t, body)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.OpsGenie.Enabled {
		t.Error("expected opsgenie to be enabled")
	}
	if cfg.OpsGenie.APIKey != "secret-key" {
		t.Errorf("expected api_key %q, got %q", "secret-key", cfg.OpsGenie.APIKey)
	}
	if cfg.OpsGenie.Team != "platform" {
		t.Errorf("expected team %q, got %q", "platform", cfg.OpsGenie.Team)
	}
	if cfg.OpsGenie.Priority != "P2" {
		t.Errorf("expected priority P2, got %q", cfg.OpsGenie.Priority)
	}
}

func TestLoad_OpsGenieMissingAPIKey(t *testing.T) {
	body := `
[opsgenie]
enabled = true
`
	p := writeOpsGenieConfig(t, body)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error for missing api_key")
	}
}

func TestLoad_OpsGenieInvalidPriority(t *testing.T) {
	body := `
[opsgenie]
enabled = true
api_key = "key"
priority = "CRITICAL"
`
	p := writeOpsGenieConfig(t, body)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error for invalid priority")
	}
}
