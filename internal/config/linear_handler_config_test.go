package config

import (
	"testing"
)

func TestLinearHandlerConfig_Defaults(t *testing.T) {
	c := defaultLinearHandlerConfig()
	if c.Priority != 2 {
		t.Errorf("expected default priority 2, got %d", c.Priority)
	}
	if c.BaseURL != "https://api.linear.app" {
		t.Errorf("expected default base_url https://api.linear.app, got %s", c.BaseURL)
	}
}

func TestValidateLinearHandler_Valid(t *testing.T) {
	c := LinearHandlerConfig{
		APIKey:  "lin_api_abc",
		TeamID:  "team-1",
		BaseURL: "https://api.linear.app",
		Priority: 1,
	}
	if err := validateLinearHandler(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateLinearHandler_MissingAPIKey(t *testing.T) {
	c := LinearHandlerConfig{
		TeamID:  "team-1",
		BaseURL: "https://api.linear.app",
		Priority: 2,
	}
	err := validateLinearHandler(c)
	if err == nil {
		t.Fatal("expected error for missing api_key")
	}
}

func TestValidateLinearHandler_MissingTeamID(t *testing.T) {
	c := LinearHandlerConfig{
		APIKey:  "lin_api_abc",
		BaseURL: "https://api.linear.app",
		Priority: 2,
	}
	err := validateLinearHandler(c)
	if err == nil {
		t.Fatal("expected error for missing team_id")
	}
}

func TestValidateLinearHandler_InvalidPriority(t *testing.T) {
	c := LinearHandlerConfig{
		APIKey:   "lin_api_abc",
		TeamID:   "team-1",
		BaseURL:  "https://api.linear.app",
		Priority: 9,
	}
	err := validateLinearHandler(c)
	if err == nil {
		t.Fatal("expected error for out-of-range priority")
	}
}
