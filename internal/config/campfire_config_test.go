package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeCampfireConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.toml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_CampfireDefaults(t *testing.T) {
	cfg := defaultCampfireConfig()
	if cfg.Enabled {
		t.Error("expected enabled=false by default")
	}
	if cfg.BaseURL != "https://api.campfirenow.com" {
		t.Errorf("unexpected base_url: %s", cfg.BaseURL)
	}
}

func TestLoad_CampfireSection(t *testing.T) {
	p := writeCampfireConfig(t, `
[campfire]
enabled = true
token = "mytoken"
account_id = "12345"
room_id = "99"
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if !cfg.Campfire.Enabled {
		t.Error("expected campfire enabled")
	}
	if cfg.Campfire.Token != "mytoken" {
		t.Errorf("unexpected token: %s", cfg.Campfire.Token)
	}
	if cfg.Campfire.RoomID != "99" {
		t.Errorf("unexpected room_id: %s", cfg.Campfire.RoomID)
	}
}

func TestLoad_CampfireMissingToken(t *testing.T) {
	p := writeCampfireConfig(t, `
[campfire]
enabled = true
account_id = "12345"
room_id = "99"
`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected validation error for missing token")
	}
}

func TestLoad_CampfireMissingRoomID(t *testing.T) {
	p := writeCampfireConfig(t, `
[campfire]
enabled = true
token = "tok"
account_id = "12345"
`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected validation error for missing room_id")
	}
}
