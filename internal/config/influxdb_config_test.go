package config_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeInfluxDBConfig(t *testing.T, content string) string {
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

func TestLoad_InfluxDBDefaults(t *testing.T) {
	path := writeInfluxDBConfig(t, "[influxdb]\nenabled = false\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.InfluxDB.Measurement != "portwatch_changes" {
		t.Errorf("expected default measurement, got %q", cfg.InfluxDB.Measurement)
	}
	if cfg.InfluxDB.Timeout != 5 {
		t.Errorf("expected default timeout 5, got %d", cfg.InfluxDB.Timeout)
	}
}

func TestLoad_InfluxDBSection(t *testing.T) {
	path := writeInfluxDBConfig(t, `[influxdb]
enabled = true
url = "http://influx:8086"
token = "mytoken"
org = "myorg"
bucket = "mybucket"
measurement = "port_events"
timeout_seconds = 10
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.InfluxDB.URL != "http://influx:8086" {
		t.Errorf("unexpected URL: %q", cfg.InfluxDB.URL)
	}
	if cfg.InfluxDB.Bucket != "mybucket" {
		t.Errorf("unexpected bucket: %q", cfg.InfluxDB.Bucket)
	}
	if cfg.InfluxDB.Measurement != "port_events" {
		t.Errorf("unexpected measurement: %q", cfg.InfluxDB.Measurement)
	}
}

func TestLoad_InfluxDBMissingToken(t *testing.T) {
	path := writeInfluxDBConfig(t, `[influxdb]
enabled = true
url = "http://influx:8086"
org = "myorg"
bucket = "mybucket"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing token")
	}
}

func TestLoad_InfluxDBInvalidTimeout(t *testing.T) {
	path := writeInfluxDBConfig(t, `[influxdb]
enabled = true
url = "http://influx:8086"
token = "tok"
org = "org"
bucket = "bkt"
timeout_seconds = 0
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for timeout <= 0")
	}
}
