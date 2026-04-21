package config

import "testing"

func TestSignalWireHandlerConfig_Defaults(t *testing.T) {
	c := defaultSignalWireHandlerConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if c.ProjectID != "" || c.APIToken != "" {
		t.Error("expected empty credential fields by default")
	}
}

func TestValidateSignalWireHandler_DisabledSkipsValidation(t *testing.T) {
	c := defaultSignalWireHandlerConfig()
	if err := validateSignalWireHandler(c); err != nil {
		t.Fatalf("expected no error for disabled handler, got: %v", err)
	}
}

func TestValidateSignalWireHandler_Valid(t *testing.T) {
	c := SignalWireHandlerConfig{
		Enabled:   true,
		ProjectID: "proj-123",
		APIToken:  "token-abc",
		From:      "+15550001111",
		To:        "+15559998888",
		SpaceURL:  "example.signalwire.com",
	}
	if err := validateSignalWireHandler(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSignalWireHandler_MissingProjectID(t *testing.T) {
	c := SignalWireHandlerConfig{
		Enabled:  true,
		APIToken: "token",
		From:     "+1",
		To:       "+2",
		SpaceURL: "example.signalwire.com",
	}
	if err := validateSignalWireHandler(c); err == nil {
		t.Fatal("expected error for missing project_id")
	}
}

func TestValidateSignalWireHandler_MissingTo(t *testing.T) {
	c := SignalWireHandlerConfig{
		Enabled:   true,
		ProjectID: "proj",
		APIToken:  "token",
		From:      "+1",
		SpaceURL:  "example.signalwire.com",
	}
	if err := validateSignalWireHandler(c); err == nil {
		t.Fatal("expected error for missing to number")
	}
}

func TestSignalWireHandlerConfigFromSettings_Valid(t *testing.T) {
	s := map[string]string{
		"enabled":    "true",
		"project_id": "proj-xyz",
		"api_token":  "tok",
		"from":       "+15550000001",
		"to":         "+15550000002",
		"space_url":  "myspace.signalwire.com",
	}
	c, err := SignalWireHandlerConfigFromSettings(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ProjectID != "proj-xyz" {
		t.Errorf("expected project_id proj-xyz, got %s", c.ProjectID)
	}
	if c.SpaceURL != "myspace.signalwire.com" {
		t.Errorf("expected space_url myspace.signalwire.com, got %s", c.SpaceURL)
	}
}
