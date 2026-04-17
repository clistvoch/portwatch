package config

import (
	"fmt"
	"time"
)

// ClickHouseHandlerConfig holds resolved configuration for the ClickHouse alert handler.
type ClickHouseHandlerConfig struct {
	DSN     string
	Table   string
	Timeout time.Duration
}

func defaultClickHouseHandlerConfig() ClickHouseHandlerConfig {
	return ClickHouseHandlerConfig{
		DSN:     "",
		Table:   "portwatch_changes",
		Timeout: 5 * time.Second,
	}
}

func validateClickHouseHandler(c ClickHouseHandlerConfig) error {
	if c.DSN == "" {
		return &ValidationError{Field: "dsn", Message: "clickhouse dsn is required"}
	}
	if c.Table == "" {
		return &ValidationError{Field: "table", Message: "clickhouse table must not be empty"}
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("clickhouse timeout must be positive")
	}
	return nil
}

// ClickHouseHandlerConfigFromSettings builds a ClickHouseHandlerConfig from raw settings map.
func ClickHouseHandlerConfigFromSettings(s map[string]string) (ClickHouseHandlerConfig, error) {
	cfg := defaultClickHouseHandlerConfig()
	if v, ok := s["dsn"]; ok && v != "" {
		cfg.DSN = v
	}
	if v, ok := s["table"]; ok && v != "" {
		cfg.Table = v
	}
	if v, ok := s["timeout"]; ok && v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return cfg, fmt.Errorf("clickhouse invalid timeout %q: %w", v, err)
		}
		cfg.Timeout = d
	}
	if err := validateClickHouseHandler(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
