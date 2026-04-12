package config

import (
	"os"
	"testing"
)

func writeSNMPConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-snmp-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_SNMPDefaults(t *testing.T) {
	path := writeSNMPConfig(t, "[scan]\nstart=1\nend=1024\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SNMP.Enabled {
		t.Error("expected SNMP disabled by default")
	}
	if cfg.SNMP.Port != 162 {
		t.Errorf("expected default port 162, got %d", cfg.SNMP.Port)
	}
	if cfg.SNMP.Community != "public" {
		t.Errorf("expected default community 'public', got %q", cfg.SNMP.Community)
	}
}

func TestLoad_SNMPSection(t *testing.T) {
	path := writeSNMPConfig(t, "[scan]\nstart=1\nend=1024\n[snmp]\nenabled=true\ntarget=\"10.0.0.1\"\nport=162\ncommunity=\"private\"\nversion=\"2c\"\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.SNMP.Enabled {
		t.Error("expected SNMP enabled")
	}
	if cfg.SNMP.Target != "10.0.0.1" {
		t.Errorf("expected target '10.0.0.1', got %q", cfg.SNMP.Target)
	}
	if cfg.SNMP.Community != "private" {
		t.Errorf("expected community 'private', got %q", cfg.SNMP.Community)
	}
}

func TestLoad_SNMPMissingTarget(t *testing.T) {
	path := writeSNMPConfig(t, "[scan]\nstart=1\nend=1024\n[snmp]\nenabled=true\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing SNMP target")
	}
}

func TestLoad_SNMPInvalidVersion(t *testing.T) {
	path := writeSNMPConfig(t, "[scan]\nstart=1\nend=1024\n[snmp]\nenabled=true\ntarget=\"10.0.0.1\"\nversion=\"3\"\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for unsupported SNMP version")
	}
}
