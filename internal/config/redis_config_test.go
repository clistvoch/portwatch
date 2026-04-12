package config

import (
	"os"
	"testing"
)

func writeRedisConfig(t *testing.T, content string) string {
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

func TestLoad_RedisDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Redis.Enabled {
		t.Error("expected redis disabled by default")
	}
	if cfg.Redis.Address != "localhost:6379" {
		t.Errorf("unexpected default address: %s", cfg.Redis.Address)
	}
	if cfg.Redis.Topic != "portwatch:changes" {
		t.Errorf("unexpected default topic: %s", cfg.Redis.Topic)
	}
	if cfg.Redis.DB != 0 {
		t.Errorf("unexpected default db: %d", cfg.Redis.DB)
	}
}

func TestLoad_RedisSection(t *testing.T) {
	path := writeRedisConfig(t, `
[redis]
enabled = true
address = "redis.example.com:6379"
password = "secret"
db = 2
topic = "alerts"
tls_enable = true
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Redis.Enabled {
		t.Error("expected redis enabled")
	}
	if cfg.Redis.Address != "redis.example.com:6379" {
		t.Errorf("unexpected address: %s", cfg.Redis.Address)
	}
	if cfg.Redis.DB != 2 {
		t.Errorf("unexpected db: %d", cfg.Redis.DB)
	}
	if !cfg.Redis.TLSEnable {
		t.Error("expected tls_enable true")
	}
}

func TestLoad_RedisMissingAddress(t *testing.T) {
	path := writeRedisConfig(t, `
[redis]
enabled = true
address = ""
topic = "portwatch:changes"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateRedis(cfg.Redis); err == nil {
		t.Error("expected error for missing address")
	}
}

func TestLoad_RedisInvalidDB(t *testing.T) {
	path := writeRedisConfig(t, `
[redis]
enabled = true
address = "localhost:6379"
topic = "portwatch:changes"
db = 20
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateRedis(cfg.Redis); err == nil {
		t.Error("expected error for invalid db value")
	}
}
