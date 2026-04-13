package config

import "fmt"

// ClickHouseConfig holds configuration for the ClickHouse alert handler.
type ClickHouseConfig struct {
	Enabled  bool   `toml:"enabled"`
	DSN      string `toml:"dsn"`
	Database string `toml:"database"`
	Table    string `toml:"table"`
	Timeout  int    `toml:"timeout_seconds"`
}

func defaultClickHouseConfig() ClickHouseConfig {
	return ClickHouseConfig{
		Enabled:  false,
		DSN:      "",
		Database: "portwatch",
		Table:    "port_changes",
		Timeout:  5,
	}
}

func validateClickHouse(c ClickHouseConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.DSN == "" {
		return &ValidationError{Field: "clickhouse.dsn", Msg: "DSN is required"}
	}
	if c.Database == "" {
		return &ValidationError{Field: "clickhouse.database", Msg: "database name is required"}
	}
	if c.Table == "" {
		return &ValidationError{Field: "clickhouse.table", Msg: "table name is required"}
	}
	if c.Timeout <= 0 {
		return &ValidationError{Field: "clickhouse.timeout_seconds", Msg: fmt.Sprintf("timeout must be positive, got %d", c.Timeout)}
	}
	return nil
}
