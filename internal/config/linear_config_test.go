package config

import (
	"os"
	"testing"
)

func writeLinearConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_LinearDefaults(t *testing.T) {
	cfg := defaultLinearConfig()
	if cfg.Enabled {
		t.Error("expected enabled to be false by default")
	}
	if cfg.BaseURL != "https://api.linear.app" {
		t.Errorf("unexpected default base_url: %s", cfg.BaseURL)
	}
	if cfg.Priority != 0 {
		t.Errorf("unexpected default priority: %d", cfg.Priority)
	}
}

func TestLoad_LinearSection(t *testing.T) {
	path := writeLinearConfig(t, `
[linear]
enabled = true
api_key = "lin_api_abc123"
team_id = "TEAM-1"
project_id = "PROJ-42"
priority = 2
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !cfg.Linear.Enabled {
		t.Error("expected linear to be enabled")
	}
	if cfg.Linear.APIKey != "lin_api_abc123" {
		t.Errorf("unexpected api_key: %s", cfg.Linear.APIKey)
	}
	if cfg.Linear.TeamID != "TEAM-1" {
		t.Errorf("unexpected team_id: %s", cfg.Linear.TeamID)
	}
	if cfg.Linear.Priority != 2 {
		t.Errorf("unexpected priority: %d", cfg.Linear.Priority)
	}
}

func TestLoad_LinearMissingAPIKey(t *testing.T) {
	path := writeLinearConfig(t, `
[linear]
enabled = true
team_id = "TEAM-1"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing api_key")
	}
}

func TestLoad_LinearInvalidPriority(t *testing.T) {
	path := writeLinearConfig(t, `
[linear]
enabled = true
api_key = "lin_api_abc123"
team_id = "TEAM-1"
priority = 9
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for invalid priority")
	}
}
