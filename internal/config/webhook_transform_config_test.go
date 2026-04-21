package config

import (
	"os"
	"testing"

	"github.com/BurntSushi/toml"
)

func writeWebhookTransformConfig(t *testing.T, content string) string {
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

func TestLoad_WebhookTransformDefaults(t *testing.T) {
	c := defaultWebhookTransformConfig()
	if c.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if c.ContentType != "application/json" {
		t.Errorf("unexpected default ContentType: %s", c.ContentType)
	}
	if !c.IncludeHost {
		t.Error("expected IncludeHost=true by default")
	}
}

func TestLoad_WebhookTransformSection(t *testing.T) {
	path := writeWebhookTransformConfig(t, `
[webhook_transform]
enabled = true
template = "{{.Port}} {{.State}}"
content_type = "text/plain"
include_host = false
`)
	var raw struct {
		WebhookTransform WebhookTransformConfig `toml:"webhook_transform"`
	}
	if _, err := toml.DecodeFile(path, &raw); err != nil {
		t.Fatal(err)
	}
	if !raw.WebhookTransform.Enabled {
		t.Error("expected Enabled=true")
	}
	if raw.WebhookTransform.Template != "{{.Port}} {{.State}}" {
		t.Errorf("unexpected Template: %s", raw.WebhookTransform.Template)
	}
	if raw.WebhookTransform.ContentType != "text/plain" {
		t.Errorf("unexpected ContentType: %s", raw.WebhookTransform.ContentType)
	}
}

func TestLoad_WebhookTransformInvalidContentType(t *testing.T) {
	c := WebhookTransformConfig{
		Enabled:     true,
		ContentType: "application/xml",
	}
	if err := validateWebhookTransform(c); err == nil {
		t.Error("expected error for unsupported content_type")
	}
}

func TestLoad_WebhookTransformDisabledSkipsValidation(t *testing.T) {
	c := WebhookTransformConfig{
		Enabled:     false,
		ContentType: "",
	}
	if err := validateWebhookTransform(c); err != nil {
		t.Errorf("expected no error when disabled, got: %v", err)
	}
}
