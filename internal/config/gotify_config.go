package config

import "fmt"

// GotifyConfig holds configuration for the Gotify push notification handler.
type GotifyConfig struct {
	Enabled  bool   `toml:"enabled"`
	URL      string `toml:"url"`
	Token    string `toml:"token"`
	Priority int    `toml:"priority"`
}

func defaultGotifyConfig() GotifyConfig {
	return GotifyConfig{
		Enabled:  false,
		URL:      "",
		Token:    "",
		Priority: 5,
	}
}

func validateGotify(c GotifyConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return ValidationError{Field: "gotify.url", Message: "url is required when gotify is enabled"}
	}
	if c.Token == "" {
		return ValidationError{Field: "gotify.token", Message: "token is required when gotify is enabled"}
	}
	if c.Priority < 0 || c.Priority > 10 {
		return ValidationError{Field: "gotify.priority", Message: fmt.Sprintf("priority must be between 0 and 10, got %d", c.Priority)}
	}
	return nil
}
