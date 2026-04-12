package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeScriptConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.toml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeScriptConfig: %v", err)
	}
	return p
}

func TestLoad_ScriptDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Script.Enabled {
		t.Error("expected script disabled by default")
	}
	if cfg.Script.Timeout != 10 {
		t.Errorf("expected default timeout 10, got %d", cfg.Script.Timeout)
	}
}

func TestLoad_ScriptSection(t *testing.T) {
	p := writeScriptConfig(t, `
[script]
enabled = true
path = "/usr/local/bin/alert.sh"
timeout_seconds = 30
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.Script.Enabled {
		t.Error("expected script enabled")
	}
	if cfg.Script.Path != "/usr/local/bin/alert.sh" {
		t.Errorf("unexpected path: %s", cfg.Script.Path)
	}
	if cfg.Script.Timeout != 30 {
		t.Errorf("expected timeout 30, got %d", cfg.Script.Timeout)
	}
}

func TestLoad_ScriptMissingPath(t *testing.T) {
	p := writeScriptConfig(t, `
[script]
enabled = true
`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for missing script path")
	}
}

func TestLoad_ScriptInvalidTimeout(t *testing.T) {
	p := writeScriptConfig(t, `
[script]
enabled = true
path = "/usr/local/bin/alert.sh"
timeout_seconds = 0
`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for invalid timeout")
	}
}
