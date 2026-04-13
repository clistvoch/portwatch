package config

import "fmt"

// HealthCheckConfig holds configuration for the HTTP health check endpoint.
type HealthCheckConfig struct {
	Enabled bool   `toml:"enabled"`
	Address string `toml:"address"`
	Path    string `toml:"path"`
}

func defaultHealthCheckConfig() HealthCheckConfig {
	return HealthCheckConfig{
		Enabled: false,
		Address: ":9110",
		Path:    "/healthz",
	}
}

func validateHealthCheck(c HealthCheckConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Address == "" {
		return fmt.Errorf("healthcheck: address must not be empty")
	}
	if len(c.Path) == 0 || c.Path[0] != '/' {
		return fmt.Errorf("healthcheck: path must start with '/'")
	}
	return nil
}
