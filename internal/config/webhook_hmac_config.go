package config

import "fmt"

// WebhookHMACConfig holds settings for HMAC-signed webhook delivery.
type WebhookHMACConfig struct {
	Enabled   bool   `toml:"enabled"`
	Secret    string `toml:"secret"`
	Algorithm string `toml:"algorithm"` // sha256 or sha512
	Header    string `toml:"header"`
}

func defaultWebhookHMACConfig() WebhookHMACConfig {
	return WebhookHMACConfig{
		Enabled:   false,
		Secret:    "",
		Algorithm: "sha256",
		Header:    "X-Portwatch-Signature",
	}
}

func validateWebhookHMAC(c WebhookHMACConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Secret == "" {
		return &ValidationError{Field: "hmac.secret", Message: "secret is required when HMAC is enabled"}
	}
	switch c.Algorithm {
	case "sha256", "sha512":
		// valid
	default:
		return &ValidationError{
			Field:   "hmac.algorithm",
			Message: fmt.Sprintf("unsupported algorithm %q: must be sha256 or sha512", c.Algorithm),
		}
	}
	if c.Header == "" {
		return &ValidationError{Field: "hmac.header", Message: "header name must not be empty"}
	}
	return nil
}
