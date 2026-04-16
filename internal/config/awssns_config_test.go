package config

import (
	"os"
	"testing"
)

func writeAWSSNSConfig(t *testing.T, content string) string {
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

func TestLoad_AWSSNSDefaults(t *testing.T) {
	path := writeAWSSNSConfig(t, "[general]\ninterval = 60\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AWSSNS.Enabled {
		t.Error("expected awssns disabled by default")
	}
	if cfg.AWSSNS.Region != "us-east-1" {
		t.Errorf("expected default region us-east-1, got %s", cfg.AWSSNS.Region)
	}
	if cfg.AWSSNS.Subject != "portwatch alert" {
		t.Errorf("unexpected default subject: %s", cfg.AWSSNS.Subject)
	}
}

func TestLoad_AWSSNSSection(t *testing.T) {
	path := writeAWSSNSConfig(t, `
[awssns]
enabled = true
region = "eu-west-1"
topic_arn = "arn:aws:sns:eu-west-1:123456789012:portwatch"
access_key = "AKIAIOSFODNN7EXAMPLE"
secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
subject = "Port Alert"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.AWSSNS.Enabled {
		t.Error("expected awssns enabled")
	}
	if cfg.AWSSNS.Region != "eu-west-1" {
		t.Errorf("unexpected region: %s", cfg.AWSSNS.Region)
	}
	if cfg.AWSSNS.TopicARN != "arn:aws:sns:eu-west-1:123456789012:portwatch" {
		t.Errorf("unexpected topic_arn: %s", cfg.AWSSNS.TopicARN)
	}
}

func TestLoad_AWSSNSMissingTopicARN(t *testing.T) {
	path := writeAWSSNSConfig(t, `
[awssns]
enabled = true
region = "us-east-1"
access_key = "key"
secret_key = "secret"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing topic_arn")
	}
}

func TestLoad_AWSSNSMissingAccessKey(t *testing.T) {
	path := writeAWSSNSConfig(t, `
[awssns]
enabled = true
region = "us-east-1"
topic_arn = "arn:aws:sns:us-east-1:000000000000:test"
secret_key = "secret"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing access_key")
	}
}
