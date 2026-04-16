package config

import "fmt"

type AWSLambdaConfig struct {
	Enabled      bool   `toml:"enabled"`
	FunctionName string `toml:"function_name"`
	Region       string `toml:"region"`
	AccessKey    string `toml:"access_key"`
	SecretKey    string `toml:"secret_key"`
	InvocationType string `toml:"invocation_type"`
}

func defaultAWSLambdaConfig() AWSLambdaConfig {
	return AWSLambdaConfig{
		Enabled:        false,
		Region:         "us-east-1",
		InvocationType: "Event",
	}
}

func validateAWSLambda(c AWSLambdaConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.FunctionName == "" {
		return &ValidationError{Field: "aws_lambda.function_name", Message: "function name is required"}
	}
	if c.AccessKey == "" {
		return &ValidationError{Field: "aws_lambda.access_key", Message: "access key is required"}
	}
	if c.SecretKey == "" {
		return &ValidationError{Field: "aws_lambda.secret_key", Message: "secret key is required"}
	}
	valid := map[string]bool{"Event": true, "RequestResponse": true, "DryRun": true}
	if !valid[c.InvocationType] {
		return &ValidationError{Field: "aws_lambda.invocation_type", Message: fmt.Sprintf("invalid invocation type %q", c.InvocationType)}
	}
	return nil
}
