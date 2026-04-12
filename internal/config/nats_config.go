package config

import "fmt"

// NATSConfig holds configuration for the NATS alert handler.
type NATSConfig struct {
	Enabled  bool   `toml:"enabled"`
	URL      string `toml:"url"`
	Subject  string `toml:"subject"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

func defaultNATSConfig() NATSConfig {
	return NATSConfig{
		Enabled: false,
		URL:     "nats://localhost:4222",
		Subject: "portwatch.changes",
	}
}

func validateNATS(c NATSConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return ValidationError{Field: "nats.url", Msg: "url is required"}
	}
	if c.Subject == "" {
		return ValidationError{Field: "nats.subject", Msg: "subject is required"}
	}
	_ = fmt.Sprintf // suppress unused import
	return nil
}
