package config

import "fmt"

type XMPPConfig struct {
	Enabled  bool   `toml:"enabled"`
	Server   string `toml:"server"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	To       string `toml:"to"`
	UseTLS   bool   `toml:"use_tls"`
}

func defaultXMPPConfig() XMPPConfig {
	return XMPPConfig{
		Enabled: false,
		Port:    5222,
		UseTLS:  true,
	}
}

func validateXMPP(c XMPPConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Server == "" {
		return fmt.Errorf("xmpp: server is required")
	}
	if c.Username == "" {
		return fmt.Errorf("xmpp: username is required")
	}
	if c.Password == "" {
		return fmt.Errorf("xmpp: password is required")
	}
	if c.To == "" {
		return fmt.Errorf("xmpp: to (recipient JID) is required")
	}
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("xmpp: port must be between 1 and 65535")
	}
	return nil
}
