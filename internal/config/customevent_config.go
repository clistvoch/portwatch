package config

import "fmt"

// CustomEventConfig holds settings for the custom event webhook handler.
type CustomEventConfig struct {
	Enabled   bool   `toml:"enabled"`
	URL       string `toml:"url"`
	Method    string `toml:"method"`
	Secret    string `toml:"secret"`
	TimeoutMs int    `toml:"timeout_ms"`
}

func defaultCustomEventConfig() CustomEventConfig {
	return CustomEventConfig{
		Enabled:   false,
		URL:       "",
		Method:    "POST",
		Secret:    "",
		TimeoutMs: 5000,
	}
}

func validateCustomEvent(cfg CustomEventConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.URL == "" {
		return &ValidationError{Field: "custom_event.url", Message: "url is required when custom_event is enabled"}
	}
	if cfg.Method != "POST" && cfg.Method != "PUT" && cfg.Method != "PATCH" {
		return &ValidationError{Field: "custom_event.method", Message: fmt.Sprintf("unsupported method %q: must be POST, PUT, or PATCH", cfg.Method)}
	}
	if cfg.TimeoutMs <= 0 {
		return &ValidationError{Field: "custom_event.timeout_ms", Message: "timeout_ms must be greater than 0"}
	}
	return nil
}
