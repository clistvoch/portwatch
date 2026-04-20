package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSignalWireConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.toml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_SignalWireDefaults(t *testing.T) {
	cfg := defaultSignalWireConfig()
	if cfg.SpaceURL != "https://example.signalwire.com" {
		t.Errorf("unexpected default space_url: %s", cfg.SpaceURL)
	}
	if cfg.ProjectID != "" {
		t.Errorf("expected empty project_id, got %s", cfg.ProjectID)
	}
}

func TestLoad_SignalWireSection(t *testing.T) {
	p := writeSignalWireConfig(t, `
[signalwire]
project_id = "proj-123"
api_token  = "tok-abc"
space_url  = "https://myspace.signalwire.com"
from       = "+15550001111"
to         = "+15559998888"
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	sw := cfg.SignalWire
	if sw.ProjectID != "proj-123" {
		t.Errorf("expected proj-123, got %s", sw.ProjectID)
	}
	if sw.APIToken != "tok-abc" {
		t.Errorf("expected tok-abc, got %s", sw.APIToken)
	}
	if sw.From != "+15550001111" {
		t.Errorf("unexpected from: %s", sw.From)
	}
}

func TestLoad_SignalWireMissingProjectID(t *testing.T) {
	cfg := SignalWireConfig{
		APIToken: "tok",
		SpaceURL: "https://x.signalwire.com",
		From:     "+1",
		To:       "+2",
	}
	if err := validateSignalWire(cfg); err == nil {
		t.Error("expected error for missing project_id")
	}
}

func TestLoad_SignalWireMissingTo(t *testing.T) {
	cfg := SignalWireConfig{
		ProjectID: "proj",
		APIToken:  "tok",
		SpaceURL:  "https://x.signalwire.com",
		From:      "+1",
	}
	if err := validateSignalWire(cfg); err == nil {
		t.Error("expected error for missing to number")
	}
}
