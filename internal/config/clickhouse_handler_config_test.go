package config

import (
	"testing"
	"time"
)

func TestClickHouseHandlerConfig_Defaults(t *testing.T) {
	cfg := defaultClickHouseHandlerConfig()
	if cfg.Table != "portwatch_changes" {
		t.Errorf("expected default table 'portwatch_changes', got %q", cfg.Table)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", cfg.Timeout)
	}
}

func TestValidateClickHouseHandler_Valid(t *testing.T) {
	cfg := ClickHouseHandlerConfig{
		DSN:     "clickhouse://localhost:9000/default",
		Table:   "portwatch_changes",
		Timeout: 5 * time.Second,
	}
	if err := validateClickHouseHandler(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateClickHouseHandler_MissingDSN(t *testing.T) {
	cfg := defaultClickHouseHandlerConfig()
	if err := validateClickHouseHandler(cfg); err == nil {
		t.Error("expected error for missing dsn")
	}
}

func TestValidateClickHouseHandler_EmptyTable(t *testing.T) {
	cfg := ClickHouseHandlerConfig{
		DSN:     "clickhouse://localhost:9000/default",
		Table:   "",
		Timeout: 5 * time.Second,
	}
	if err := validateClickHouseHandler(cfg); err == nil {
		t.Error("expected error for empty table")
	}
}

func TestClickHouseHandlerConfigFromSettings_Valid(t *testing.T) {
	s := map[string]string{
		"dsn":     "clickhouse://localhost:9000/default",
		"table":   "events",
		"timeout": "10s",
	}
	cfg, err := ClickHouseHandlerConfigFromSettings(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Table != "events" {
		t.Errorf("expected table 'events', got %q", cfg.Table)
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("expected timeout 10s, got %v", cfg.Timeout)
	}
}

func TestClickHouseHandlerConfigFromSettings_InvalidTimeout(t *testing.T) {
	s := map[string]string{
		"dsn":     "clickhouse://localhost:9000/default",
		"timeout": "not-a-duration",
	}
	_, err := ClickHouseHandlerConfigFromSettings(s)
	if err == nil {
		t.Error("expected error for invalid timeout")
	}
}
