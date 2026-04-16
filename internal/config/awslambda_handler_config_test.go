package config

import (
	"testing"
)

func TestAWSLambdaHandlerConfig_Defaults(t *testing.T) {
	c := defaultAWSLambdaHandlerConfig()
	if c.Region != "us-east-1" {
		t.Errorf("expected region us-east-1, got %s", c.Region)
	}
	if c.InvocationType != "Event" {
		t.Errorf("expected invocation_type Event, got %s", c.InvocationType)
	}
	if c.Timeout != 10 {
		t.Errorf("expected timeout 10, got %d", c.Timeout)
	}
}

func TestValidateAWSLambdaHandler_Valid(t *testing.T) {
	c := defaultAWSLambdaHandlerConfig()
	c.FunctionName = "my-fn"
	c.AccessKey = "AKID"
	c.SecretKey = "secret"
	if err := validateAWSLambdaHandler(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateAWSLambdaHandler_MissingFunctionName(t *testing.T) {
	c := defaultAWSLambdaHandlerConfig()
	c.AccessKey = "AKID"
	c.SecretKey = "secret"
	if err := validateAWSLambdaHandler(c); err == nil {
		t.Fatal("expected error for missing function_name")
	}
}

func TestValidateAWSLambdaHandler_InvalidInvocationType(t *testing.T) {
	c := defaultAWSLambdaHandlerConfig()
	c.FunctionName = "fn"
	c.AccessKey = "AKID"
	c.SecretKey = "secret"
	c.InvocationType = "BadType"
	if err := validateAWSLambdaHandler(c); err == nil {
		t.Fatal("expected error for invalid invocation_type")
	}
}

func TestValidateAWSLambdaHandler_MissingKeys(t *testing.T) {
	c := defaultAWSLambdaHandlerConfig()
	c.FunctionName = "fn"
	if err := validateAWSLambdaHandler(c); err == nil {
		t.Fatal("expected error for missing access_key")
	}
}
