package config

import "fmt"

// CampfireHandlerConfig holds resolved settings for the Campfire alert handler.
type CampfireHandlerConfig struct {
	Token  string
	RoomID string
	APIURL string
}

func defaultCampfireHandlerConfig() CampfireHandlerConfig {
	return CampfireHandlerConfig{
		APIURL: "https://api.campfirenow.com",
	}
}

func validateCampfireHandler(c CampfireHandlerConfig) error {
	if c.Token == "" {
		return &ValidationError{Field: "campfire.token", Msg: "token is required"}
	}
	if c.RoomID == "" {
		return &ValidationError{Field: "campfire.room_id", Msg: "room_id is required"}
	}
	if c.APIURL == "" {
		return &ValidationError{Field: "campfire.api_url", Msg: "api_url must not be empty"}
	}
	return nil
}

// CampfireHandlerConfigFromSettings builds a CampfireHandlerConfig from raw
// map values loaded from the TOML config file.
func CampfireHandlerConfigFromSettings(s map[string]string) (CampfireHandlerConfig, error) {
	cfg := defaultCampfireHandlerConfig()
	if v, ok := s["token"]; ok && v != "" {
		cfg.Token = v
	}
	if v, ok := s["room_id"]; ok && v != "" {
		cfg.RoomID = v
	}
	if v, ok := s["api_url"]; ok && v != "" {
		cfg.APIURL = v
	}
	if err := validateCampfireHandler(cfg); err != nil {
		return cfg, fmt.Errorf("campfire handler: %w", err)
	}
	return cfg, nil
}
