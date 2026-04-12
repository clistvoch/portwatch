package config

// ScriptConfig holds configuration for the script alert handler.
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

func validateScript(cfg ScriptConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.Path == "" {
		return &ValidationError{Field: "script.path", Message: "path is required when script handler is enabled"}
	}
	if cfg.Timeout <= 0 {
		return &ValidationError{Field: "script.timeout_seconds", Message: "must be greater than 0"}
	}
	return nil
}
