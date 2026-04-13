package config

import (
	"os"
	"testing"
)

func writeHealthCheckConfig(t *testing.T, body string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(body); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_HealthCheckDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.HealthCheck.Enabled {
		t.Error("expected healthcheck disabled by default")
	}
	if cfg.HealthCheck.Address != ":9110" {
		t.Errorf("unexpected default address: %s", cfg.HealthCheck.Address)
	}
	if cfg.HealthCheck.Path != "/healthz" {
		t.Errorf("unexpected default path: %s", cfg.HealthCheck.Path)
	}
}

func TestLoad_HealthCheckSection(t *testing.T) {
	path := writeHealthCheckConfig(t, `
[healthcheck]
enabled = true
address = ":8080"
path = "/health"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.HealthCheck.Enabled {
		t.Error("expected healthcheck enabled")
	}
	if cfg.HealthCheck.Address != ":8080" {
		t.Errorf("unexpected address: %s", cfg.HealthCheck.Address)
	}
	if cfg.HealthCheck.Path != "/health" {
		t.Errorf("unexpected path: %s", cfg.HealthCheck.Path)
	}
}

func TestLoad_HealthCheckInvalidPath(t *testing.T) {
	path := writeHealthCheckConfig(t, `
[healthcheck]
enabled = true
address = ":9110"
path = "healthz"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for path missing leading slash")
	}
}

func TestLoad_HealthCheckMissingAddress(t *testing.T) {
	path := writeHealthCheckConfig(t, `
[healthcheck]
enabled = true
address = ""
path = "/healthz"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}
