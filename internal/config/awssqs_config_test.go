package config

import (
	"os"
	"testing"
)

func writeAWSSQSConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.toml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestLoad_AWSSQSDefaults(t *testing.T) {
	cfg := defaultAWSSQSConfig()
	if cfg.Region != "us-east-1" {
		t.Errorf("expected us-east-1, got %s", cfg.Region)
	}
	if cfg.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.MessageGroupID != "portwatch" {
		t.Errorf("expected portwatch, got %s", cfg.MessageGroupID)
	}
}

func TestLoad_AWSSQSSection(t *testing.T) {
	path := writeAWSSQSConfig(t, `
[awssqs]
enabled = true
queue_url = "https://sqs.us-east-1.amazonaws.com/123/myqueue"
region = "us-west-2"
access_key = "AKID"
secret_key = "SECRET"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.AWSSQS.Enabled {
		t.Error("expected enabled")
	}
	if cfg.AWSSQS.Region != "us-west-2" {
		t.Errorf("unexpected region: %s", cfg.AWSSQS.Region)
	}
}

func TestLoad_AWSSQSMissingQueueURL(t *testing.T) {
	c := AWSSQSConfig{Enabled: true, AccessKey: "k", SecretKey: "s"}
	if err := validateAWSSQS(c); err == nil {
		t.Error("expected error for missing queue_url")
	}
}

func TestLoad_AWSSQSMissingAccessKey(t *testing.T) {
	c := AWSSQSConfig{Enabled: true, QueueURL: "https://sqs.example.com/q", SecretKey: "s"}
	if err := validateAWSSQS(c); err == nil {
		t.Error("expected error for missing access_key")
	}
}
