package config

import "fmt"

type LinearHandlerConfig struct {
	APIKey      string
	TeamID      string
	ProjectID   string
	Priority    int
	LabelIDs    []string
	BaseURL     string
	AssigneeID  string
}

func defaultLinearHandlerConfig() LinearHandlerConfig {
	return LinearHandlerConfig{
		Priority: 2,
		BaseURL:  "https://api.linear.app",
	}
}

func validateLinearHandler(c LinearHandlerConfig) error {
	if c.APIKey == "" {
		return &ValidationError{Field: "linear.api_key", Msg: "api_key is required"}
	}
	if c.TeamID == "" {
		return &ValidationError{Field: "linear.team_id", Msg: "team_id is required"}
	}
	if c.Priority < 0 || c.Priority > 4 {
		return &ValidationError{Field: "linear.priority", Msg: fmt.Sprintf("priority must be 0-4, got %d", c.Priority)}
	}
	if c.BaseURL == "" {
		return &ValidationError{Field: "linear.base_url", Msg: "base_url must not be empty"}
	}
	return nil
}
