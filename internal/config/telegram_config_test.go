package config

import (
	"os"
	"testing"
)

func writeTelegramConfig(t *testing.T, body string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(body)
	_ = f.Close()
	return f.Name()
}

func TestLoad_TelegramDefaults(t *testing.T) {
	path := writeTelegramConfig(t, "")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Telegram.Enabled {
		t.Error("expected telegram disabled by default")
	}
	if cfg.Telegram.ParseMode != "HTML" {
		t.Errorf("expected default parse_mode HTML, got %q", cfg.Telegram.ParseMode)
	}
}

func TestLoad_TelegramSection(t *testing.T) {
	path := writeTelegramConfig(t, `
[telegram]
enabled = true
bot_token = "123:ABC"
chat_id = "-100999"
parse_mode = "Markdown"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Telegram.Enabled {
		t.Error("expected telegram enabled")
	}
	if cfg.Telegram.BotToken != "123:ABC" {
		t.Errorf("unexpected bot_token: %s", cfg.Telegram.BotToken)
	}
	if cfg.Telegram.ChatID != "-100999" {
		t.Errorf("unexpected chat_id: %s", cfg.Telegram.ChatID)
	}
	if cfg.Telegram.ParseMode != "Markdown" {
		t.Errorf("unexpected parse_mode: %s", cfg.Telegram.ParseMode)
	}
}

func TestLoad_TelegramMissingBotToken(t *testing.T) {
	path := writeTelegramConfig(t, `
[telegram]
enabled = true
chat_id = "-100999"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing bot_token")
	}
}

func TestLoad_TelegramInvalidParseMode(t *testing.T) {
	path := writeTelegramConfig(t, `
[telegram]
enabled = true
bot_token = "tok"
chat_id = "123"
parse_mode = "plain"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for invalid parse_mode")
	}
}
