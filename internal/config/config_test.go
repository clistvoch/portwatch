package config

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatalf("encode config: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.StartPort != 1 || cfg.EndPort != 1024 {
		t.Errorf("unexpected default ports: %d-%d", cfg.StartPort, cfg.EndPort)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("unexpected default interval: %v", cfg.Interval)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	path := writeTempConfig(t, map[string]any{
		"start_port": 8000,
		"end_port":   9000,
		"interval":   60000000000, // 60s in nanoseconds
	})
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.StartPort != 8000 || cfg.EndPort != 9000 {
		t.Errorf("ports not loaded: %d-%d", cfg.StartPort, cfg.EndPort)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidate_InvalidRange(t *testing.T) {
	cfg := DefaultConfig()
	cfg.StartPort = 900
	cfg.EndPort = 100
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for inverted port range")
	}
}

func TestValidate_ShortInterval(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Interval = 500 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for interval < 1s")
	}
}
