package config

import (
	"os"
	"testing"
)

func writeAppriseConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestLoad_AppriseDefaults(t *testing.T) {
	cfg := defaultAppriseConfig()
	if cfg.Enabled {
		t.Error("expected enabled=false")
	}
	if cfg.Tag != "portwatch" {
		t.Errorf("unexpected tag: %s", cfg.Tag)
	}
	if cfg.Title != "PortWatch Alert" {
		t.Errorf("unexpected title: %s", cfg.Title)
	}
}

func TestLoad_AppriseSection(t *testing.T) {
	path := writeAppriseConfig(t, `
[apprise]
enabled = true
url = "http://localhost:8000/notify"
tag = "alerts"
title = "My Alert"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Apprise.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.Apprise.URL != "http://localhost:8000/notify" {
		t.Errorf("unexpected url: %s", cfg.Apprise.URL)
	}
	if cfg.Apprise.Tag != "alerts" {
		t.Errorf("unexpected tag: %s", cfg.Apprise.Tag)
	}
}

func TestLoad_AppriseMissingURL(t *testing.T) {
	c := AppriseConfig{Enabled: true, URL: ""}
	if err := validateApprise(c); err == nil {
		t.Error("expected error for missing url")
	}
}
