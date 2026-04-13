package config

import (
	"os"
	"testing"
)

func writeHipChatConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_HipChatDefaults(t *testing.T) {
	path := writeHipChatConfig(t, "[scan]\nports = \"80-81\"\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HipChat.Enabled {
		t.Error("expected HipChat disabled by default")
	}
	if cfg.HipChat.BaseURL != "https://api.hipchat.com" {
		t.Errorf("unexpected base_url: %s", cfg.HipChat.BaseURL)
	}
	if cfg.HipChat.Color != "yellow" {
		t.Errorf("unexpected default color: %s", cfg.HipChat.Color)
	}
}

func TestLoad_HipChatSection(t *testing.T) {
	path := writeHipChatConfig(t, `
[scan]
ports = "80-81"

[hipchat]
enabled = true
room_id = "42"
auth_token = "secret"
color = "green"
notify = true
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.HipChat.Enabled {
		t.Error("expected HipChat enabled")
	}
	if cfg.HipChat.RoomID != "42" {
		t.Errorf("unexpected room_id: %s", cfg.HipChat.RoomID)
	}
	if cfg.HipChat.Color != "green" {
		t.Errorf("unexpected color: %s", cfg.HipChat.Color)
	}
	if !cfg.HipChat.Notify {
		t.Error("expected notify=true")
	}
}

func TestLoad_HipChatMissingRoomID(t *testing.T) {
	path := writeHipChatConfig(t, `
[scan]
ports = "80-81"

[hipchat]
enabled = true
auth_token = "secret"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing room_id")
	}
}

func TestLoad_HipChatInvalidColor(t *testing.T) {
	path := writeHipChatConfig(t, `
[scan]
ports = "80-81"

[hipchat]
enabled = true
room_id = "42"
auth_token = "secret"
color = "blue"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid color")
	}
}
