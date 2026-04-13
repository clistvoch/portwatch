package config

import "fmt"

type StatusPageConfig struct {
	Enabled     bool   `toml:"enabled"`
	APIKey      string `toml:"api_key"`
	PageID      string `toml:"page_id"`
	ComponentID string `toml:"component_id"`
	BaseURL     string `toml:"base_url"`
}

func defaultStatusPageHandlerConfig() StatusPageConfig {
	return StatusPageConfig{
		Enabled: false,
		BaseURL: "https://api.statuspage.io/v1",
	}
}

func validateStatusPageHandler(c StatusPageConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return &ValidationError{Field: "statuspage.api_key", Msg: "api_key is required"}
	}
	if c.PageID == "" {
		return &ValidationError{Field: "statuspage.page_id", Msg: "page_id is required"}
	}
	if c.ComponentID == "" {
		return &ValidationError{Field: "statuspage.component_id", Msg: "component_id is required"}
	}
	if c.BaseURL == "" {
		return fmt.Errorf("statuspage.base_url must not be empty")
	}
	return nil
}
