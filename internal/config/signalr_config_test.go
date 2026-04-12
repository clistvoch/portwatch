package config

import (
	"os"
	"testing"
)

func writeSignalRConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_SignalRDefaults(t *testing.T) {
	path := writeSignalRConfig(t, "[general]\ninterval_sec = 60\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.SignalR.Enabled {
		t.Error("expected SignalR disabled by default")
	}
	if cfg.SignalR.Hub != "portwatch" {
		t.Errorf("expected hub 'portwatch', got %q", cfg.SignalR.Hub)
	}
	if cfg.SignalR.TimeoutSec != 10 {
		t.Errorf("expected timeout 10, got %d", cfg.SignalR.TimeoutSec)
	}
}

func TestLoad_SignalRSection(t *testing.T) {
	path := writeSignalRConfig(t, `
[general]
interval_sec = 60

[signalr]
enabled = true
endpoint = "https://example.service.signalr.net"
access_key = "secret123"
hub = "alerts"
timeout_sec = 5
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.SignalR.Enabled {
		t.Error("expected SignalR enabled")
	}
	if cfg.SignalR.Endpoint != "https://example.service.signalr.net" {
		t.Errorf("unexpected endpoint: %q", cfg.SignalR.Endpoint)
	}
	if cfg.SignalR.Hub != "alerts" {
		t.Errorf("unexpected hub: %q", cfg.SignalR.Hub)
	}
	if cfg.SignalR.TimeoutSec != 5 {
		t.Errorf("unexpected timeout: %d", cfg.SignalR.TimeoutSec)
	}
}

func TestLoad_SignalRMissingEndpoint(t *testing.T) {
	path := writeSignalRConfig(t, `
[general]
interval_sec = 60

[signalr]
enabled = true
access_key = "secret123"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
}
