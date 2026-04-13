package config_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeClickHouseConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestLoad_ClickHouseDefaults(t *testing.T) {
	path := writeClickHouseConfig(t, "")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ClickHouse.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.ClickHouse.Database != "portwatch" {
		t.Errorf("expected database=portwatch, got %q", cfg.ClickHouse.Database)
	}
	if cfg.ClickHouse.Table != "port_changes" {
		t.Errorf("expected table=port_changes, got %q", cfg.ClickHouse.Table)
	}
	if cfg.ClickHouse.Timeout != 5 {
		t.Errorf("expected timeout=5, got %d", cfg.ClickHouse.Timeout)
	}
}

func TestLoad_ClickHouseSection(t *testing.T) {
	path := writeClickHouseConfig(t, `
[clickhouse]
enabled = true
dsn = "clickhouse://localhost:9000"
database = "monitoring"
table = "events"
timeout_seconds = 10
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.ClickHouse.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.ClickHouse.DSN != "clickhouse://localhost:9000" {
		t.Errorf("unexpected DSN: %q", cfg.ClickHouse.DSN)
	}
	if cfg.ClickHouse.Database != "monitoring" {
		t.Errorf("unexpected database: %q", cfg.ClickHouse.Database)
	}
	if cfg.ClickHouse.Timeout != 10 {
		t.Errorf("unexpected timeout: %d", cfg.ClickHouse.Timeout)
	}
}

func TestLoad_ClickHouseMissingDSN(t *testing.T) {
	path := writeClickHouseConfig(t, `
[clickhouse]
enabled = true
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing DSN")
	}
}

func TestLoad_ClickHouseInvalidTimeout(t *testing.T) {
	path := writeClickHouseConfig(t, `
[clickhouse]
enabled = true
dsn = "clickhouse://localhost:9000"
timeout_seconds = -1
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for negative timeout")
	}
}
