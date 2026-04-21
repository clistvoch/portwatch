package config

import (
	"os"
	"testing"
)

func writeWebhookHMACConfig(t *testing.T, content string) string {
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

func TestLoad_WebhookHMACDefaults(t *testing.T) {
	cfg := defaultWebhookHMACConfig()
	if cfg.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.Algorithm != "sha256" {
		t.Errorf("expected algorithm=sha256, got %s", cfg.Algorithm)
	}
	if cfg.Header != "X-Portwatch-Signature" {
		t.Errorf("unexpected default header: %s", cfg.Header)
	}
}

func TestLoad_WebhookHMACSection(t *testing.T) {
	path := writeWebhookHMACConfig(t, `
[webhook_hmac]
enabled = true
secret = "mysecret"
algorithm = "sha512"
header = "X-Sig"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.WebhookHMAC.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.WebhookHMAC.Secret != "mysecret" {
		t.Errorf("unexpected secret: %s", cfg.WebhookHMAC.Secret)
	}
	if cfg.WebhookHMAC.Algorithm != "sha512" {
		t.Errorf("unexpected algorithm: %s", cfg.WebhookHMAC.Algorithm)
	}
}

func TestLoad_WebhookHMACMissingSecret(t *testing.T) {
	cfg := WebhookHMACConfig{Enabled: true, Algorithm: "sha256", Header: "X-Sig"}
	if err := validateWebhookHMAC(cfg); err == nil {
		t.Error("expected error for missing secret")
	}
}

func TestLoad_WebhookHMACInvalidAlgorithm(t *testing.T) {
	cfg := WebhookHMACConfig{Enabled: true, Secret: "s", Algorithm: "md5", Header: "X-Sig"}
	if err := validateWebhookHMAC(cfg); err == nil {
		t.Error("expected error for invalid algorithm")
	}
}
