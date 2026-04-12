package config

import (
	"os"
	"testing"
)

func writeElasticsearchConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.toml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestLoad_ElasticsearchDefaults(t *testing.T) {
	path := writeElasticsearchConfig(t, "")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Elasticsearch.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.Elasticsearch.URL != "http://localhost:9200" {
		t.Errorf("unexpected default URL: %s", cfg.Elasticsearch.URL)
	}
	if cfg.Elasticsearch.Index != "portwatch" {
		t.Errorf("unexpected default index: %s", cfg.Elasticsearch.Index)
	}
	if cfg.Elasticsearch.TimeoutSec != 5 {
		t.Errorf("unexpected default timeout: %d", cfg.Elasticsearch.TimeoutSec)
	}
}

func TestLoad_ElasticsearchSection(t *testing.T) {
	path := writeElasticsearchConfig(t, `
[elasticsearch]
enabled = true
url = "http://es.example.com:9200"
index = "alerts"
username = "admin"
password = "secret"
timeout_sec = 10
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Elasticsearch.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.Elasticsearch.URL != "http://es.example.com:9200" {
		t.Errorf("unexpected URL: %s", cfg.Elasticsearch.URL)
	}
	if cfg.Elasticsearch.Index != "alerts" {
		t.Errorf("unexpected index: %s", cfg.Elasticsearch.Index)
	}
	if cfg.Elasticsearch.TimeoutSec != 10 {
		t.Errorf("unexpected timeout: %d", cfg.Elasticsearch.TimeoutSec)
	}
}

func TestLoad_ElasticsearchMissingURL(t *testing.T) {
	path := writeElasticsearchConfig(t, `
[elasticsearch]
enabled = true
url = ""
index = "portwatch"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}
