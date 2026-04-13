package config

import "fmt"

// StatusPageConfig holds configuration for Atlassian Statuspage alerting.
type StatusPageConfig struct {
	Enabled   bool   `toml:"enabled"`
	APIKey    string `toml:"api_key"`
	PageID    string `toml:"page_id"`
	BaseURL   string `toml:"base_url"`
	ComponentID string `toml:"component_id"`
}

func defaultStatusPageConfig() StatusPageConfig {
	return StatusPageConfig{
		Enabled:   false,
		APIKey:    "",
		PageID:    "",
		BaseURL:   "https://api.statuspage.io/v1",
		ComponentID: "",
	}
}

func validateStatusPage(cfg StatusPageConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.APIKey == "" {
		return fmt.Errorf("statuspage: api_key is required")
	}
	if cfg.PageID == "" {
		return fmt.Errorf("statuspage: page_id is required")
	}
	if cfg.ComponentID == "" {
		return fmt.Errorf("statuspage: component_id is required")
	}
	if cfg.BaseURL == "" {
		return fmt.Errorf("statuspage: base_url must not be empty")
	}
	return nil
}
