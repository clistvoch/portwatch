package config

import "fmt"

// InfluxDBConfig holds configuration for the InfluxDB alert handler.
type InfluxDBConfig struct {
	Enabled     bool   `toml:"enabled"`
	URL         string `toml:"url"`
	Token       string `toml:"token"`
	Org         string `toml:"org"`
	Bucket      string `toml:"bucket"`
	Measurement string `toml:"measurement"`
	Timeout     int    `toml:"timeout_seconds"`
}

func defaultInfluxDBConfig() InfluxDBConfig {
	return InfluxDBConfig{
		Enabled:     false,
		URL:         "http://localhost:8086",
		Measurement: "portwatch_changes",
		Timeout:     5,
	}
}

func validateInfluxDB(cfg InfluxDBConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.URL == "" {
		return ValidationError{Field: "influxdb.url", Msg: "url is required"}
	}
	if cfg.Token == "" {
		return ValidationError{Field: "influxdb.token", Msg: "token is required"}
	}
	if cfg.Org == "" {
		return ValidationError{Field: "influxdb.org", Msg: "org is required"}
	}
	if cfg.Bucket == "" {
		return ValidationError{Field: "influxdb.bucket", Msg: "bucket is required"}
	}
	if cfg.Timeout <= 0 {
		return ValidationError{Field: "influxdb.timeout_seconds", Msg: fmt.Sprintf("must be > 0, got %d", cfg.Timeout)}
	}
	return nil
}
