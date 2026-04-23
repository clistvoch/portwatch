package config

import "fmt"

type IRCConfig struct {
	Enabled  bool   `toml:"enabled"`
	Server   string `toml:"server"`
	Port     int    `toml:"port"`
	Nick     string `toml:"nick"`
	Channel  string `toml:"channel"`
	Password string `toml:"password"`
	TLS      bool   `toml:"tls"`
}

func defaultIRCConfig() IRCConfig {
	return IRCConfig{
		Enabled: false,
		Server:  "irc.libera.chat",
		Port:    6667,
		Nick:    "portwatch",
		Channel: "#alerts",
		TLS:     false,
	}
}

func validateIRC(c IRCConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Server == "" {
		return &ValidationError{Field: "irc.server", Message: "server address is required"}
	}
	if c.Port <= 0 || c.Port > 65535 {
		return &ValidationError{Field: "irc.port", Message: fmt.Sprintf("invalid port %d", c.Port)}
	}
	if c.Nick == "" {
		return &ValidationError{Field: "irc.nick", Message: "nick is required"}
	}
	if c.Channel == "" {
		return &ValidationError{Field: "irc.channel", Message: "channel is required"}
	}
	return nil
}
