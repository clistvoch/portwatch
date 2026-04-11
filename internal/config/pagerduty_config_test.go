package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writePagerDutyConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.toml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

func TestLoad_PagerDutyDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.PagerDuty.Enabled {
		t.Error("expected pagerduty disabled by default")
	}
	if cfg.PagerDuty.RoutingKey != "" {
		t.Errorf("expected empty routing_key, got %q", cfg.PagerDuty.RoutingKey)
	}
}

func TestLoad_PagerDutySection(t *testing.T) {
	path := writePagerDutyConfig(t, `
[pagerduty]
enabled = true
routing_key = "abc123"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.PagerDuty.Enabled {
		t.Error("expected pagerduty enabled")
	}
	if cfg.PagerDuty.RoutingKey != "abc123" {
		t.Errorf("routing_key = %q, want abc123", cfg.PagerDuty.RoutingKey)
	}
}

func TestLoad_PagerDutyMissingRoutingKey(t *testing.T) {
	path := writePagerDutyConfig(t, `
[pagerduty]
enabled = true
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := Validate(cfg); err == nil {
		t.Error("expected validation error for missing routing_key")
	}
}
