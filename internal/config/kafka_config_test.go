package config

import (
	"os"
	"testing"
)

func writeKafkaConfig(t *testing.T, content string) string {
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

func TestLoad_KafkaDefaults(t *testing.T) {
	path := writeKafkaConfig(t, "[general]\ninterval = 30\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Kafka.Enabled {
		t.Error("expected kafka disabled by default")
	}
	if cfg.Kafka.Topic != "portwatch-alerts" {
		t.Errorf("expected default topic, got %q", cfg.Kafka.Topic)
	}
	if cfg.Kafka.ClientID != "portwatch" {
		t.Errorf("expected default client_id, got %q", cfg.Kafka.ClientID)
	}
}

func TestLoad_KafkaSection(t *testing.T) {
	path := writeKafkaConfig(t, `
[general]
interval = 30

[kafka]
enabled = true
brokers = ["localhost:9092", "localhost:9093"]
topic = "my-topic"
client_id = "my-client"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Kafka.Enabled {
		t.Error("expected kafka enabled")
	}
	if len(cfg.Kafka.Brokers) != 2 {
		t.Errorf("expected 2 brokers, got %d", len(cfg.Kafka.Brokers))
	}
	if cfg.Kafka.Topic != "my-topic" {
		t.Errorf("expected my-topic, got %q", cfg.Kafka.Topic)
	}
}

func TestLoad_KafkaMissingBrokers(t *testing.T) {
	path := writeKafkaConfig(t, `
[general]
interval = 30

[kafka]
enabled = true
topic = "portwatch-alerts"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing brokers")
	}
}
