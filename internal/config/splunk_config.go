package config

import "fmt"

// SplunkConfig holds configuration for the Splunk HEC alert handler.
type SplunkConfig struct {
	Enabled    bool   `toml:"enabled"`
	URL        string `toml:"url"`
	Token      string `toml:"token"`
	Index      string `toml:"index"`
	SourceType string `toml:"source_type"`
	Timeout    int    `toml:"timeout_seconds"`
}

func defaultSplunkConfig() SplunkConfig {
	return SplunkConfig{
		Enabled:    false,
		URL:        "",
		Token:      "",
		Index:      "main",
		SourceType: "portwatch",
		Timeout:    10,
	}
}

func validateSplunk(c SplunkConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return &ValidationError{Field: "splunk.url", Msg: "url is required when splunk is enabled"}
	}
	if c.Token == "" {
		return &ValidationError{Field: "splunk.token", Msg: "token is required when splunk is enabled"}
	}
	if c.Timeout <= 0 {
		return &ValidationError{Field: "splunk.timeout_seconds", Msg: fmt.Sprintf("timeout must be positive, got %d", c.Timeout)}
	}
	return nil
}
