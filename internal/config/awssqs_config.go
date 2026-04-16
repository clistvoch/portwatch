package config

import "fmt"

type AWSSQSConfig struct {
	Enabled         bool   `toml:"enabled"`
	QueueURL        string `toml:"queue_url"`
	Region          string `toml:"region"`
	AccessKey       string `toml:"access_key"`
	SecretKey       string `toml:"secret_key"`
	MessageGroupID  string `toml:"message_group_id"`
}

func defaultAWSSQSConfig() AWSSQSConfig {
	return AWSSQSConfig{
		Enabled:        false,
		Region:         "us-east-1",
		MessageGroupID: "portwatch",
	}
}

func validateAWSSQS(c AWSSQSConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.QueueURL == "" {
		return fmt.Errorf("awssqs: queue_url is required")
	}
	if c.AccessKey == "" {
		return fmt.Errorf("awssqs: access_key is required")
	}
	if c.SecretKey == "" {
		return fmt.Errorf("awssqs: secret_key is required")
	}
	return nil
}
