package config

import (
	"testing"

	"github.com/BurntSushi/toml"
)

func writeXMPPConfig(t *testing.T, content string) string {
	t.Helper()
	return writeTempConfig(t, content)
}

func TestLoad_XMPPDefaults(t *testing.T) {
	cfg := defaultXMPPConfig()
	if cfg.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.Port != 5222 {
		t.Errorf("expected port 5222, got %d", cfg.Port)
	}
	if !cfg.UseTLS {
		t.Error("expected use_tls=true by default")
	}
}

func TestLoad_XMPPSection(t *testing.T) {
	path := writeXMPPConfig(t, `
[xmpp]
enabled = true
server = "jabber.example.com"
port = 5222
username = "bot@example.com"
password = "secret"
to = "admin@example.com"
use_tls = true
`)
	var wrapper struct {
		XMPP XMPPConfig `toml:"xmpp"`
	}
	if _, err := toml.DecodeFile(path, &wrapper); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if wrapper.XMPP.Server != "jabber.example.com" {
		t.Errorf("unexpected server: %s", wrapper.XMPP.Server)
	}
	if wrapper.XMPP.To != "admin@example.com" {
		t.Errorf("unexpected to: %s", wrapper.XMPP.To)
	}
}

func TestLoad_XMPPMissingServer(t *testing.T) {
	cfg := XMPPConfig{
		Enabled:  true,
		Username: "bot@example.com",
		Password: "secret",
		To:       "admin@example.com",
		Port:     5222,
	}
	if err := validateXMPP(cfg); err == nil {
		t.Error("expected error for missing server")
	}
}

func TestLoad_XMPPMissingTo(t *testing.T) {
	cfg := XMPPConfig{
		Enabled:  true,
		Server:   "jabber.example.com",
		Username: "bot@example.com",
		Password: "secret",
		Port:     5222,
	}
	if err := validateXMPP(cfg); err == nil {
		t.Error("expected error for missing to")
	}
}
