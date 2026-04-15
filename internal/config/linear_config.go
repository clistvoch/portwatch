package config

import "fmt"

type LinearConfig struct {
	Enabled   bool   `toml:"enabled"`
	APIKey    string `toml:"api_key"`
	TeamID    string `toml:"team_id"`
	ProjectID string `toml:"project_id"`
	BaseURL   string `toml:"base_url"`
	Priority  int    `toml:"priority"`
}

func defaultLinearConfig() LinearConfig {
	return LinearConfig{
		Enabled:  false,
		BaseURL:  "https://api.linear.app",
		Priority: 0,
	}
}

func validateLinear(cfg LinearConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.APIKey == "" {
		return &ValidationError{Field: "linear.api_key", Message: "api_key is required"}
	}
	if cfg.TeamID == "" {
		return &ValidationError{Field: "linear.team_id", Message: "team_id is required"}
	}
	if cfg.BaseURL == "" {
		return &ValidationError{Field: "linear.base_url", Message: "base_url must not be empty"}
	}
	if cfg.Priority < 0 || cfg.Priority > 4 {
		return &ValidationError{Field: "linear.priority", Message: fmt.Sprintf("priority must be between 0 and 4, got %d", cfg.Priority)}
	}
	return nil
}
