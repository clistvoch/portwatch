package config

import "fmt"

// SignalWireHandlerConfig holds resolved configuration for the SignalWire handler.
type SignalWireHandlerConfig struct {
	Enabled     bool
	ProjectID   string
	APIToken    string
	From        string
	To          string
	SpaceURL    string
}

func defaultSignalWireHandlerConfig() SignalWireHandlerConfig {
	return SignalWireHandlerConfig{
		Enabled: false,
	}
}

func validateSignalWireHandler(c SignalWireHandlerConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.ProjectID == "" {
		return fmt.Errorf("signalwire: project_id is required")
	}
	if c.APIToken == "" {
		return fmt.Errorf("signalwire: api_token is required")
	}
	if c.From == "" {
		return fmt.Errorf("signalwire: from number is required")
	}
	if c.To == "" {
		return fmt.Errorf("signalwire: to number is required")
	}
	if c.SpaceURL == "" {
		return fmt.Errorf("signalwire: space_url is required")
	}
	return nil
}

// SignalWireHandlerConfigFromSettings builds a SignalWireHandlerConfig from
// the raw settings map loaded by the config file parser.
func SignalWireHandlerConfigFromSettings(s map[string]string) (SignalWireHandlerConfig, error) {
	c := defaultSignalWireHandlerConfig()
	if v, ok := s["enabled"]; ok {
		c.Enabled = v == "true"
	}
	if v, ok := s["project_id"]; ok {
		c.ProjectID = v
	}
	if v, ok := s["api_token"]; ok {
		c.APIToken = v
	}
	if v, ok := s["from"]; ok {
		c.From = v
	}
	if v, ok := s["to"]; ok {
		c.To = v
	}
	if v, ok := s["space_url"]; ok {
		c.SpaceURL = v
	}
	if err := validateSignalWireHandler(c); err != nil {
		return c, err
	}
	return c, nil
}
