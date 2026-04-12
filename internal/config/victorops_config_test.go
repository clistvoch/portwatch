package config

import (
	"os"
	"testing"
)

func writeVictorOpsConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_VictorOpsDefaults(t *testing.T) {
	path := writeVictorOpsConfig(t, "[general]\ninterval = 30\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VictorOps.Enabled {
		t.Error("expected victorops disabled by default")
	}
	if cfg.VictorOps.MessageType != "CRITICAL" {
		t.Errorf("expected default message_type CRITICAL, got %q", cfg.VictorOps.MessageType)
	}
}

func TestLoad_VictorOpsSection(t *testing.T) {
	path := writeVictorOpsConfig(t, `
[victorops]
enabled = true
webhook_url = "https://alert.victorops.com/integrations/generic/123/alert/token"
routing_key = "ops-team"
message_type = "WARNING"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.VictorOps.Enabled {
		t.Error("expected victorops enabled")
	}
	if cfg.VictorOps.RoutingKey != "ops-team" {
		t.Errorf("expected routing_key ops-team, got %q", cfg.VictorOps.RoutingKey)
	}
	if cfg.VictorOps.MessageType != "WARNING" {
		t.Errorf("expected message_type WARNING, got %q", cfg.VictorOps.MessageType)
	}
}

func TestLoad_VictorOpsMissingWebhook(t *testing.T) {
	path := writeVictorOpsConfig(t, `
[victorops]
enabled = true
routing_key = "ops-team"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing webhook_url")
	}
}

func TestLoad_VictorOpsInvalidMessageType(t *testing.T) {
	path := writeVictorOpsConfig(t, `
[victorops]
enabled = true
webhook_url = "https://alert.victorops.com/integrations/generic/123/alert/token"
routing_key = "ops-team"
message_type = "UNKNOWN"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid message_type")
	}
}
