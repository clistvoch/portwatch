package config

import "fmt"

type AWSLambdaHandlerConfig struct {
	FunctionName   string `toml:"function_name"`
	Region         string `toml:"region"`
	AccessKey      string `toml:"access_key"`
	SecretKey      string `toml:"secret_key"`
	InvocationType string `toml:"invocation_type"`
	Timeout        int    `toml:"timeout_seconds"`
}

func defaultAWSLambdaHandlerConfig() AWSLambdaHandlerConfig {
	return AWSLambdaHandlerConfig{
		Region:         "us-east-1",
		InvocationType: "Event",
		Timeout:        10,
	}
}

func validateAWSLambdaHandler(c AWSLambdaHandlerConfig) error {
	if c.FunctionName == "" {
		return fmt.Errorf("awslambda: function_name is required")
	}
	if c.AccessKey == "" {
		return fmt.Errorf("awslambda: access_key is required")
	}
	if c.SecretKey == "" {
		return fmt.Errorf("awslambda: secret_key is required")
	}
	if c.InvocationType != "Event" && c.InvocationType != "RequestResponse" {
		return fmt.Errorf("awslambda: invocation_type must be Event or RequestResponse")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("awslambda: timeout_seconds must be positive")
	}
	return nil
}
