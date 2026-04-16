package config_test

import (
	"os"
	"testing"

	"github.com/example/portwatch/internal/config"
)

func writeAWSLambdaConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.toml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestLoad_AWSLambdaDefaults(t *testing.T) {
	path := writeAWSLambdaConfig(t, "")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AWSLambda.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.AWSLambda.Region != "us-east-1" {
		t.Errorf("unexpected region: %s", cfg.AWSLambda.Region)
	}
	if cfg.AWSLambda.InvocationType != "Event" {
		t.Errorf("unexpected invocation_type: %s", cfg.AWSLambda.InvocationType)
	}
}

func TestLoad_AWSLambdaSection(t *testing.T) {
	path := writeAWSLambdaConfig(t, `
[aws_lambda]
enabled = true
function_name = "my-func"
region = "eu-west-1"
access_key = "AKID"
secret_key = "SECRET"
invocation_type = "RequestResponse"
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.AWSLambda.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.AWSLambda.FunctionName != "my-func" {
		t.Errorf("unexpected function_name: %s", cfg.AWSLambda.FunctionName)
	}
	if cfg.AWSLambda.InvocationType != "RequestResponse" {
		t.Errorf("unexpected invocation_type: %s", cfg.AWSLambda.InvocationType)
	}
}

func TestLoad_AWSLambdaMissingFunctionName(t *testing.T) {
	path := writeAWSLambdaConfig(t, `
[aws_lambda]
enabled = true
access_key = "AKID"
secret_key = "SECRET"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing function_name")
	}
}

func TestLoad_AWSLambdaInvalidInvocationType(t *testing.T) {
	path := writeAWSLambdaConfig(t, `
[aws_lambda]
enabled = true
function_name = "fn"
access_key = "AKID"
secret_key = "SECRET"
invocation_type = "BadType"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid invocation_type")
	}
}
