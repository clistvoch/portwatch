package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	StartPort  int           `json:"start_port"`
	EndPort    int           `json:"end_port"`
	Interval   time.Duration `json:"interval"`
	LogFile    string        `json:"log_file"`
	LogPrefix  string        `json:"log_prefix"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		StartPort: 1,
		EndPort:   1024,
		Interval:  30 * time.Second,
		LogFile:   "",
		LogPrefix: "[portwatch]",
	}
}

// Load reads a JSON config file from path and merges it over defaults.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, cfg.Validate()
}

// Validate checks that the Config fields are within acceptable bounds.
func (c *Config) Validate() error {
	if c.StartPort < 1 || c.StartPort > 65535 {
		return errors.New("config: start_port must be between 1 and 65535")
	}
	if c.EndPort < 1 || c.EndPort > 65535 {
		return errors.New("config: end_port must be between 1 and 65535")
	}
	if c.StartPort > c.EndPort {
		return errors.New("config: start_port must not be greater than end_port")
	}
	if c.Interval < time.Second {
		return errors.New("config: interval must be at least 1s")
	}
	return nil
}
