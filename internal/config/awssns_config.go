package config

import "fmt"

type AWSSNSConfig struct {
	Enabled   bool   `toml:"enabled"`
	Region    string `toml:"region"`
	TopicARN  string `toml:"topic_arn"`
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
	Subject   string `toml:"subject"`
}

func defaultAWSSNSConfig() AWSSNSConfig {
	return AWSSNSConfig{
		Enabled:  false,
		Region:   "us-east-1",
		Subject:  "portwatch alert",
	}
}

func validateAWSSNS(c AWSSNSConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.TopicARN == "" {
		return &ValidationError{Field: "awssns.topic_arn", Message: "topic ARN is required"}
	}
	if c.Region == "" {
		return &ValidationError{Field: "awssns.region", Message: "region is required"}
	}
	if c.AccessKey == "" {
		return &ValidationError{Field: "awssns.access_key", Message: "access key is required"}
	}
	if c.SecretKey == "" {
		return fmt.Errorf("awssns: secret_key is required")
	}
	return nil
}
