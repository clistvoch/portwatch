package config

import "fmt"

// SyslogConfig holds configuration for the syslog alert handler.
type SyslogConfig struct {
	Enabled  bool   `toml:"enabled"`
	Network  string `toml:"network"`  // "" for local, "tcp", "udp"
	Address  string `toml:"address"`  // host:port for remote syslog
	Tag      string `toml:"tag"`
	Priority string `toml:"priority"` // "info", "warning", "err"
}

func defaultSyslogConfig() SyslogConfig {
	return SyslogConfig{
		Enabled:  false,
		Network:  "",
		Address:  "",
		Tag:      "portwatch",
		Priority: "info",
	}
}

func validateSyslog(c SyslogConfig) error {
	if !c.Enabled {
		return nil
	}
	validPriorities := map[string]bool{
		"emerg": true, "alert": true, "crit": true, "err": true,
		"warning": true, "notice": true, "info": true, "debug": true,
	}
	if !validPriorities[c.Priority] {
		return fmt.Errorf("syslog: invalid priority %q", c.Priority)
	}
	if c.Network != "" && c.Network != "tcp" && c.Network != "udp" {
		return fmt.Errorf("syslog: invalid network %q, must be \"tcp\", \"udp\", or empty for local", c.Network)
	}
	if c.Network != "" && c.Address == "" {
		return fmt.Errorf("syslog: address is required when network is set")
	}
	if c.Tag == "" {
		return fmt.Errorf("syslog: tag must not be empty")
	}
	return nil
}
