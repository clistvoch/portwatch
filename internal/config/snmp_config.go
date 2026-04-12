package config

import "fmt"

// SNMPConfig holds configuration for SNMP trap alerting.
type SNMPConfig struct {
	Enabled   bool   `toml:"enabled"`
	Target    string `toml:"target"`
	Port      int    `toml:"port"`
	Community string `toml:"community"`
	Version   string `toml:"version"`
}

func defaultSNMPConfig() SNMPConfig {
	return SNMPConfig{
		Enabled:   false,
		Target:    "",
		Port:      162,
		Community: "public",
		Version:   "2c",
	}
}

func validateSNMP(c SNMPConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Target == "" {
		return ValidationError{Field: "snmp.target", Msg: "target host is required"}
	}
	if c.Port < 1 || c.Port > 65535 {
		return ValidationError{Field: "snmp.port", Msg: fmt.Sprintf("invalid port %d", c.Port)}
	}
	if c.Version != "1" && c.Version != "2c" {
		return ValidationError{Field: "snmp.version", Msg: fmt.Sprintf("unsupported version %q, use '1' or '2c'", c.Version)}
	}
	return nil
}
