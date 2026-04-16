package config_test

import (
	"testing"

	"github.com/patrickward/portwatch/internal/config"
)

func writeGCPPubSubConfig(t *testing.T, content string) string {
	t.Helper()
	return writeTempConfig(t, content)
}

func TestLoad_GCPPubSubDefaults(t *testing.T) {
	path := writeGCPPubSubConfig(t, "")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.GCPPubSub.Enabled {
		t.Error("expected enabled=false")
	}
	if cfg.GCPPubSub.TopicID != "portwatch-alerts" {
		t.Errorf("unexpected topic_id: %s", cfg.GCPPubSub.TopicID)
	}
}

func TestLoad_GCPPubSubSection(t *testing.T) {
	path := writeGCPPubSubConfig(t, `
[gcppubsub]
enabled = true
project_id = "my-project"
topic_id = "my-topic"
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.GCPPubSub.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.GCPPubSub.ProjectID != "my-project" {
		t.Errorf("unexpected project_id: %s", cfg.GCPPubSub.ProjectID)
	}
}

func TestLoad_GCPPubSubMissingProjectID(t *testing.T) {
	path := writeGCPPubSubConfig(t, `
[gcppubsub]
enabled = true
topic_id = "my-topic"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error")
	}
}
