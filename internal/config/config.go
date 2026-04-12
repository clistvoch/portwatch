package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config holds the full portwatch configuration.
type Config struct {
	PortRange string `toml:"port_range"`
	Interval  int    `toml:"interval_seconds"`
	StateFile string `toml:"state_file"`
	LogFile   string `toml:"log_file"`
	Webhook   string `toml:"webhook_url"`
	Email     Email  `toml:"email"`
}

// Email holds SMTP alert settings.
type Email struct {
	Enabled  bool     `toml:"enabled"`
	Host     string   `toml:"host"`
	Port     int      `toml:"port"`
	Username string   `toml:"username"`
	Password string   `toml:"password"`
	From     string   `toml:"from"`
	To       []string `toml:"to"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		PortRange: "1-1024",
		Interval:  60,
		StateFile: "/tmp/portwatch.state",
	}
}

// Load reads a TOML config file into a Config struct.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return cfg, fmt.Errorf("config file not found: %s", path)
	}
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}
	return cfg, Validate(cfg)
}

// Validate checks Config fields for correctness.
func Validate(cfg Config) error {
	if cfg.Interval <= 0 {
		return fmt.Errorf("interval_seconds must be > 0")
	}
	if err := validatePortRange(cfg.PortRange); err != nil {
		return err
	}
	if cfg.Email.Enabled {
		if cfg.Email.Host == "" {
			return fmt.Errorf("email.host is required when email is enabled")
		}
		if len(cfg.Email.To) == 0 {
			return fmt.Errorf("email.to must have at least one recipient when email is enabled")
		}
	}
	return nil
}

// validatePortRange checks that the port range string is well-formed and within
// the valid port number bounds [1, 65535] with lo <= hi.
func validatePortRange(portRange string) error {
	var lo, hi int
	if _, err := fmt.Sscanf(portRange, "%d-%d", &lo, &hi); err != nil {
		return fmt.Errorf("invalid port_range %q", portRange)
	}
	if lo < 1 || hi > 65535 || lo > hi {
		return fmt.Errorf("port_range out of bounds: %s", portRange)
	}
	return nil
}
