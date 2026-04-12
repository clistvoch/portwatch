package config_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeNATSConfig(t *testing.T, body string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(body); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_NATSDefaults(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.NATS.Enabled {
		t.Error("expected NATS disabled by default")
	}
	if cfg.NATS.URL != "nats://localhost:4222" {
		t.Errorf("unexpected default URL: %s", cfg.NATS.URL)
	}
	if cfg.NATS.Subject != "portwatch.changes" {
		t.Errorf("unexpected default subject: %s", cfg.NATS.Subject)
	}
}

func TestLoad_NATSSection(t *testing.T) {
	path := writeNATSConfig(t, `
[nats]
enabled = true
url = "nats://nats.example.com:4222"
subject = "alerts.ports"
username = "user"
password = "secret"
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.NATS.Enabled {
		t.Error("expected NATS enabled")
	}
	if cfg.NATS.URL != "nats://nats.example.com:4222" {
		t.Errorf("unexpected URL: %s", cfg.NATS.URL)
	}
	if cfg.NATS.Subject != "alerts.ports" {
		t.Errorf("unexpected subject: %s", cfg.NATS.Subject)
	}
	if cfg.NATS.Username != "user" {
		t.Errorf("unexpected username: %s", cfg.NATS.Username)
	}
}

func TestLoad_NATSMissingURL(t *testing.T) {
	path := writeNATSConfig(t, `
[nats]
enabled = true
url = ""
subject = "alerts.ports"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing NATS URL")
	}
}

func TestLoad_NATSMissingSubject(t *testing.T) {
	path := writeNATSConfig(t, `
[nats]
enabled = true
url = "nats://localhost:4222"
subject = ""
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty NATS subject")
	}
}
