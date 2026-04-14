package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeCloudWatchConfig(t *testing.T, content string) string {
	t.Helper()
	return writeTempConfig(t, content)
}

func TestLoad_CloudWatchDefaults(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CloudWatch.Enabled {
		t.Error("expected cloudwatch disabled by default")
	}
	if cfg.CloudWatch.Region != "us-east-1" {
		t.Errorf("expected region us-east-1, got %q", cfg.CloudWatch.Region)
	}
	if cfg.CloudWatch.Namespace != "PortWatch" {
		t.Errorf("expected namespace PortWatch, got %q", cfg.CloudWatch.Namespace)
	}
}

func TestLoad_CloudWatchSection(t *testing.T) {
	path := writeCloudWatchConfig(t, `
[cloudwatch]
enabled = true
region = "eu-west-1"
namespace = "MyApp"
metric_name = "PortEvent"
access_key = "AKIAIOSFODNN7EXAMPLE"
secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.CloudWatch.Enabled {
		t.Error("expected cloudwatch enabled")
	}
	if cfg.CloudWatch.Region != "eu-west-1" {
		t.Errorf("unexpected region: %q", cfg.CloudWatch.Region)
	}
	if cfg.CloudWatch.MetricName != "PortEvent" {
		t.Errorf("unexpected metric_name: %q", cfg.CloudWatch.MetricName)
	}
}

func TestLoad_CloudWatchMissingAccessKey(t *testing.T) {
	path := writeCloudWatchConfig(t, `
[cloudwatch]
enabled = true
region = "us-east-1"
namespace = "PortWatch"
metric_name = "PortChange"
secret_key = "somesecret"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing access_key")
	}
}

func TestLoad_CloudWatchMissingSecretKey(t *testing.T) {
	path := writeCloudWatchConfig(t, `
[cloudwatch]
enabled = true
region = "us-east-1"
namespace = "PortWatch"
metric_name = "PortChange"
access_key = "somekey"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing secret_key")
	}
}
