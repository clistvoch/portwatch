package config

import (
	"testing"
	"time"
)

func TestResolve_Defaults(t *testing.T) {
	cfg, err := Resolve(&Flags{})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if cfg.StartPort != 1 || cfg.EndPort != 1024 {
		t.Errorf("expected default ports 1-1024, got %d-%d", cfg.StartPort, cfg.EndPort)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
}

func TestResolve_FlagOverrides(t *testing.T) {
	cfg, err := Resolve(&Flags{
		StartPort: 3000,
		EndPort:   4000,
		Interval:  5 * time.Second,
		LogFile:   "/tmp/portwatch.log",
	})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if cfg.StartPort != 3000 || cfg.EndPort != 4000 {
		t.Errorf("flag override failed: ports %d-%d", cfg.StartPort, cfg.EndPort)
	}
	if cfg.Interval != 5*time.Second {
		t.Errorf("flag override failed: interval %v", cfg.Interval)
	}
	if cfg.LogFile != "/tmp/portwatch.log" {
		t.Errorf("flag override failed: log file %q", cfg.LogFile)
	}
}

func TestResolve_FileAndFlagMerge(t *testing.T) {
	path := writeTempConfig(t, map[string]any{
		"start_port": 8000,
		"end_port":   9000,
		"interval":   10000000000, // 10s
	})
	cfg, err := Resolve(&Flags{
		ConfigPath: path,
		EndPort:    9500, // override only end port
	})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if cfg.StartPort != 8000 {
		t.Errorf("expected start_port 8000 from file, got %d", cfg.StartPort)
	}
	if cfg.EndPort != 9500 {
		t.Errorf("expected end_port 9500 from flag, got %d", cfg.EndPort)
	}
}

func TestResolve_InvalidFlagRange(t *testing.T) {
	_, err := Resolve(&Flags{StartPort: 5000, EndPort: 1000})
	if err == nil {
		t.Error("expected validation error for inverted range")
	}
}
