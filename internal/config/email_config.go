package config

import "fmt"

// EmailConfig holds SMTP alert configuration.
type EmailConfig struct {
	Enabled    bool   `toml:"enabled"`
	Host       string `toml:"host"`
	Port       int    `toml:"port"`
	Username   string `toml:"username"`
	Password   string `toml:"password"`
	From       string `toml:"from"`
	To         string `toml:"to"`
	Subject    string `toml:"subject"`
	SkipVerify bool   `toml:"skip_verify"`
}

func defaultEmailConfig() EmailConfig {
	return EmailConfig{
		Enabled: false,
		Port:    587,
		Subject: "portwatch: port change detected",
	}
}

func validateEmail(e EmailConfig) error {
	if !e.Enabled {
		return nil
	}
	if e.Host == "" {
		return ValidationError{Field: "email.host", Msg: "host is required when email is enabled"}
	}
	if e.Port <= 0 || e.Port > 65535 {
		return ValidationError{Field: "email.port", Msg: fmt.Sprintf("invalid port %d", e.Port)}
	}
	if e.From == "" {
		return ValidationError{Field: "email.from", Msg: "from address is required when email is enabled"}
	}
	if e.To == "" {
		return ValidationError{Field: "email.to", Msg: "to address is required when email is enabled"}
	}
	return nil
}
