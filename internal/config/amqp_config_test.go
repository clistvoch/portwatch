package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeAMQPConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.toml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

func TestLoad_AMQPDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.AMQP.Enabled {
		t.Error("expected AMQP disabled by default")
	}
	if cfg.AMQP.URL != "amqp://guest:guest@localhost:5672/" {
		t.Errorf("unexpected default URL: %s", cfg.AMQP.URL)
	}
	if cfg.AMQP.Exchange != "portwatch" {
		t.Errorf("unexpected default exchange: %s", cfg.AMQP.Exchange)
	}
	if cfg.AMQP.RoutingKey != "port.change" {
		t.Errorf("unexpected default routing key: %s", cfg.AMQP.RoutingKey)
	}
}

func TestLoad_AMQPSection(t *testing.T) {
	p := writeAMQPConfig(t, `
[amqp]
enabled = true
url = "amqp://user:pass@rabbit:5672/vhost"
exchange = "alerts"
routing_key = "portwatch.events"
durable = false
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.AMQP.Enabled {
		t.Error("expected AMQP enabled")
	}
	if cfg.AMQP.URL != "amqp://user:pass@rabbit:5672/vhost" {
		t.Errorf("unexpected URL: %s", cfg.AMQP.URL)
	}
	if cfg.AMQP.Exchange != "alerts" {
		t.Errorf("unexpected exchange: %s", cfg.AMQP.Exchange)
	}
	if cfg.AMQP.RoutingKey != "portwatch.events" {
		t.Errorf("unexpected routing key: %s", cfg.AMQP.RoutingKey)
	}
	if cfg.AMQP.Durable {
		t.Error("expected durable=false")
	}
}

func TestLoad_AMQPMissingURL(t *testing.T) {
	p := writeAMQPConfig(t, `
[amqp]
enabled = true
url = ""
exchange = "alerts"
routing_key = "portwatch.events"
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := validateAMQP(cfg.AMQP); err == nil {
		t.Error("expected validation error for missing URL")
	}
}
