package config

import (
	"fmt"
	"strings"
)

// CustomEventHandlerConfig holds resolved settings for the custom event handler.
type CustomEventHandlerConfig struct {
	URL     string
	Method  string
	Headers map[string]string
	Timeout int
}

func defaultCustomEventHandlerConfig() CustomEventHandlerConfig {
	return CustomEventHandlerConfig{
		Method:  "POST",
		Headers: map[string]string{},
		Timeout: 10,
	}
}

func validateCustomEventHandler(c CustomEventHandlerConfig) error {
	if c.URL == "" {
		return &ValidationError{Field: "url", Msg: "custom_event url is required"}
	}
	allowed := map[string]bool{"GET": true, "POST": true, "PUT": true, "PATCH": true}
	if !allowed[strings.ToUpper(c.Method)] {
		return &ValidationError{Field: "method", Msg: fmt.Sprintf("custom_event method %q is not supported", c.Method)}
	}
	if c.Timeout <= 0 {
		return &ValidationError{Field: "timeout", Msg: "custom_event timeout must be positive"}
	}
	return nil
}

// CustomEventHandlerConfigFromSettings builds a CustomEventHandlerConfig from
// the raw settings map stored in the TOML config.
func CustomEventHandlerConfigFromSettings(s map[string]string) (CustomEventHandlerConfig, error) {
	cfg := defaultCustomEventHandlerConfig()
	if v, ok := s["url"]; ok && v != "" {
		cfg.URL = v
	}
	if v, ok := s["method"]; ok && v != "" {
		cfg.Method = strings.ToUpper(v)
	}
	if v, ok := s["timeout_seconds"]; ok && v != "" {
		var t int
		if _, err := fmt.Sscanf(v, "%d", &t); err == nil {
			cfg.Timeout = t
		}
	}
	// Any key prefixed with "header_" is treated as an HTTP header.
	for k, v := range s {
		if strings.HasPrefix(k, "header_") {
			headerName := strings.TrimPrefix(k, "header_")
			cfg.Headers[headerName] = v
		}
	}
	if err := validateCustomEventHandler(cfg); err != nil {
		return CustomEventHandlerConfig{}, err
	}
	return cfg, nil
}
