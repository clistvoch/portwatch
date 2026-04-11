package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writePrometheusConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.toml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

func TestLoad_PrometheusDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Prometheus.Enabled {
		t.Error("expected prometheus disabled by default")
	}
	if cfg.Prometheus.Address != ":9090" {
		t.Errorf("expected :9090, got %s", cfg.Prometheus.Address)
	}
	if cfg.Prometheus.Path != "/metrics" {
		t.Errorf("expected /metrics, got %s", cfg.Prometheus.Path)
	}
}

func TestLoad_PrometheusSection(t *testing.T) {
	p := writePrometheusConfig(t, `
[prometheus]
enabled = true
address = ":2112"
path = "/prom"
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !cfg.Prometheus.Enabled {
		t.Error("expected prometheus enabled")
	}
	if cfg.Prometheus.Address != ":2112" {
		t.Errorf("unexpected address: %s", cfg.Prometheus.Address)
	}
	if cfg.Prometheus.Path != "/prom" {
		t.Errorf("unexpected path: %s", cfg.Prometheus.Path)
	}
}

func TestLoad_PrometheusMissingPath(t *testing.T) {
	p := writePrometheusConfig(t, `
[prometheus]
enabled = true
address = ":2112"
path = "noslash"
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if err := Validate(*cfg); err == nil {
		t.Error("expected validation error for path without leading slash")
	}
}
