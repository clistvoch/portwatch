package config

import "fmt"

// SignalRHandlerConfig holds runtime configuration for the SignalR alert handler.
type SignalRHandlerConfig struct {
	Endpoint   string
	Hub        string
	Method     string
	AccessKey  string
	TimeoutSec int
}

func defaultSignalRHandlerConfig() SignalRHandlerConfig {
	return SignalRHandlerConfig{
		Hub:        "portwatch",
		Method:     "portChange",
		TimeoutSec: 10,
	}
}

func validateSignalRHandler(c SignalRHandlerConfig) error {
	if c.Endpoint == "" {
		return &ValidationError{Field: "signalr.endpoint", Message: "endpoint is required"}
	}
	if c.Hub == "" {
		return &ValidationError{Field: "signalr.hub", Message: "hub must not be empty"}
	}
	if c.Method == "" {
		return &ValidationError{Field: "signalr.method", Message: "method must not be empty"}
	}
	if c.TimeoutSec <= 0 {
		return &ValidationError{Field: "signalr.timeout_sec", Message: fmt.Sprintf("timeout_sec must be positive, got %d", c.TimeoutSec)}
	}
	return nil
}
