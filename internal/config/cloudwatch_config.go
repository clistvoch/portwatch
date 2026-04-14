package config

import "fmt"

type CloudWatchConfig struct {
	Enabled    bool   `toml:"enabled"`
	Region     string `toml:"region"`
	Namespace  string `toml:"namespace"`
	MetricName string `toml:"metric_name"`
	AccessKey  string `toml:"access_key"`
	SecretKey  string `toml:"secret_key"`
}

func defaultCloudWatchConfig() CloudWatchConfig {
	return CloudWatchConfig{
		Enabled:    false,
		Region:     "us-east-1",
		Namespace:  "PortWatch",
		MetricName: "PortChange",
	}
}

func validateCloudWatch(c CloudWatchConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Region == "" {
		return fmt.Errorf("cloudwatch: region is required")
	}
	if c.Namespace == "" {
		return fmt.Errorf("cloudwatch: namespace is required")
	}
	if c.MetricName == "" {
		return fmt.Errorf("cloudwatch: metric_name is required")
	}
	if c.AccessKey == "" {
		return fmt.Errorf("cloudwatch: access_key is required")
	}
	if c.SecretKey == "" {
		return fmt.Errorf("cloudwatch: secret_key is required")
	}
	return nil
}
