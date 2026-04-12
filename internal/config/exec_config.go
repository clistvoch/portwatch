package config

import (
	"fmt"
	"time"
)

// ExecConfig holds configuration for the exec (script) alert handler.
type ExecConfig struct {
	Enabled bool          `toml:"enabled"`
	Path    string        `toml:"path"`
	Args    []string      `toml:"args"`
	Timeout time.Duration `toml:"timeout"`
	Shell   string        `toml:"shell"`
}

func defaultExecConfig() ExecConfig {
	return ExecConfig{
		Enabled: false,
		Path:    "",
		Args:    []string{},
		Timeout: 5 * time.Second,
		Shell:   "",
	}
}

func validateExec(c ExecConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Path == "" {
		return ValidationError{Field: "exec.path", Msg: "path is required when exec is enabled"}
	}
	if c.Timeout <= 0 {
		return ValidationError{Field: "exec.timeout", Msg: "timeout must be positive"}
	}
	if c.Timeout > 60*time.Second {
		return fmt.Errorf("exec.timeout: must not exceed 60s, got %s", c.Timeout)
	}
	return nil
}
