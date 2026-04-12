package config_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeNewRelicConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_NewRelicDefaults(t *testing.T) {
	path := writeNewRelicConfig(t, "[scan]\nstart=1\nend=1024\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.NewRelic.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.NewRelic.Region != "US" {
		t.Errorf("expected region=US, got %s", cfg.NewRelic.Region)
	}
	if cfg.NewRelic.EventType != "PortWatchAlert" {
		t.Errorf("expected event_type=PortWatchAlert, got %s", cfg.NewRelic.EventType)
	}
	if cfg.NewRelic.TimeoutSec != 10 {
		t.Errorf("expected timeout_sec=10, got %d", cfg.NewRelic.TimeoutSec)
	}
}

func TestLoad_NewRelicSection(t *testing.T) {
	path := writeNewRelicConfig(t, `
[scan]
start=1
end=1024

[newrelic]
enabled=true
api_key="NRAK-test"
account_id="123456"
region="EU"
event_type="CustomEvent"
timeout_sec=5
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.NewRelic.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.NewRelic.APIKey != "NRAK-test" {
		t.Errorf("unexpected api_key: %s", cfg.NewRelic.APIKey)
	}
	if cfg.NewRelic.Region != "EU" {
		t.Errorf("unexpected region: %s", cfg.NewRelic.Region)
	}
}

func TestLoad_NewRelicMissingAPIKey(t *testing.T) {
	path := writeNewRelicConfig(t, `
[scan]
start=1
end=1024

[newrelic]
enabled=true
account_id="123456"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing api_key")
	}
}

func TestLoad_NewRelicInvalidRegion(t *testing.T) {
	path := writeNewRelicConfig(t, `
[scan]
start=1
end=1024

[newrelic]
enabled=true
api_key="NRAK-test"
account_id="123456"
region="AP"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid region")
	}
}
