package config

import (
	"os"
	"testing"
)

func writeAzureEventHubConfig(t *testing.T, content string) string {
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

func TestLoad_AzureEventHubDefaults(t *testing.T) {
	path := writeAzureEventHubConfig(t, "")
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AzureEventHub.Enabled {
		t.Error("expected Enabled=false by default")
	}
}

func TestLoad_AzureEventHubSection(t *testing.T) {
	path := writeAzureEventHubConfig(t, `
[azure_event_hub]
enabled = true
connection_string = "Endpoint=sb://foo.servicebus.windows.net/"
event_hub_name = "portwatch"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.AzureEventHub.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.AzureEventHub.EventHubName != "portwatch" {
		t.Errorf("unexpected hub name: %s", cfg.AzureEventHub.EventHubName)
	}
}

func TestLoad_AzureEventHubMissingConnectionString(t *testing.T) {
	path := writeAzureEventHubConfig(t, `
[azure_event_hub]
enabled = true
event_hub_name = "portwatch"
`)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for missing connection_string")
	}
}

func TestLoad_AzureEventHubMissingHubName(t *testing.T) {
	path := writeAzureEventHubConfig(t, `
[azure_event_hub]
enabled = true
connection_string = "Endpoint=sb://foo.servicebus.windows.net/"
`)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for missing event_hub_name")
	}
}
