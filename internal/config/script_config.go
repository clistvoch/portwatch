package config

import "fmt"

// ScriptConfig holds configuration for the script/exec alert handler.
type ScriptConfig struct {
	Enabled bool   `toml:"enabled"`
	Path    string `toml:"path"`
	Timeout int    `toml:"timeout_seconds"`
}

func defaultScriptConfig() ScriptConfig {
	return ScriptConfig{
		Enabled: false,
		Path:    "",
		Timeout: 10,
	}
}

func validateScript(c ScriptConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Path == "" {
		return ValidationError{Field: "script.path", Msg: "path is required when script handler is enabled"}
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("script.timeout_seconds must be greater than 0")
	}
	return nil
}
