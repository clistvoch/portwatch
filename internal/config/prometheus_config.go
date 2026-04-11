package config

import "fmt"

// PrometheusConfig holds settings for the Prometheus metrics exposition handler.
type PrometheusConfig struct {
	Enabled bool   `toml:"enabled"`
	Address string `toml:"address"`
	Path    string `toml:"path"`
}

func defaultPrometheusConfig() PrometheusConfig {
	return PrometheusConfig{
		Enabled: false,
		Address: ":9090",
		Path:    "/metrics",
	}
}

func validatePrometheus(c PrometheusConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Address == "" {
		return fmt.Errorf("prometheus: address must not be empty")
	}
	if len(c.Path) == 0 || c.Path[0] != '/' {
		return fmt.Errorf("prometheus: path must start with '/'")
	}
	return nil
}
