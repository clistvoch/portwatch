package config

import (
	"os"
	"testing"
)

func writeAWSEventBridgeConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	return f.Name()
}

func TestLoad_AWSEventBridgeDefaults(t *testing.T) {
	cfg := defaultAWSEventBridgeConfig()
	if cfg.Region != "us-east-1" {
		t.Errorf("expected us-east-1, got %s", cfg.Region)
	}
	if cfg.BusName != "default" {
		t.Errorf("expected default, got %s", cfg.BusName)
	}
	if cfg.Source != "portwatch" {
		t.Errorf("expected portwatch, got %s", cfg.Source)
	}
	if cfg.Enabled {
		t.Error("expected enabled=false by default")
	}
}

func TestLoad_AWSEventBridgeSection(t *testing.T) {
	path := writeAWSEventBridgeConfig(t, `
[awseventbridge]
enabled = true
region = "eu-west-1"
bus_name = "my-bus"
source = "myapp"
detail_type = "Alert"
access_key = "AKID"
secret_key = "SECRET"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	eb := cfg.AWSEventBridge
	if !eb.Enabled {
		t.Error("expected enabled=true")
	}
	if eb.Region != "eu-west-1" {
		t.Errorf("expected eu-west-1, got %s", eb.Region)
	}
	if eb.BusName != "my-bus" {
		t.Errorf("expected my-bus, got %s", eb.BusName)
	}
	if eb.AccessKey != "AKID" {
		t.Errorf("expected AKID, got %s", eb.AccessKey)
	}
}

func TestLoad_AWSEventBridgeMissingAccessKey(t *testing.T) {
	c := defaultAWSEventBridgeConfig()
	c.Enabled = true
	c.SecretKey = "s"
	if err := validateAWSEventBridge(c); err == nil {
		t.Error("expected error for missing access_key")
	}
}

func TestLoad_AWSEventBridgeMissingSecretKey(t *testing.T) {
	c := defaultAWSEventBridgeConfig()
	c.Enabled = true
	c.AccessKey = "k"
	if err := validateAWSEventBridge(c); err == nil {
		t.Error("expected error for missing secret_key")
	}
}
