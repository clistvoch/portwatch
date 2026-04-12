package config

import "fmt"

// TelegramConfig holds configuration for Telegram bot alerting.
type TelegramConfig struct {
	Enabled  bool   `toml:"enabled"`
	BotToken string `toml:"bot_token"`
	ChatID   string `toml:"chat_id"`
	ParseMode string `toml:"parse_mode"` // HTML or Markdown
}

func defaultTelegramConfig() TelegramConfig {
	return TelegramConfig{
		Enabled:   false,
		ParseMode: "HTML",
	}
}

func validateTelegram(cfg TelegramConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.BotToken == "" {
		return &ValidationError{Field: "telegram.bot_token", Msg: "bot_token is required when telegram is enabled"}
	}
	if cfg.ChatID == "" {
		return &ValidationError{Field: "telegram.chat_id", Msg: "chat_id is required when telegram is enabled"}
	}
	if cfg.ParseMode != "HTML" && cfg.ParseMode != "Markdown" && cfg.ParseMode != "" {
		return &ValidationError{
			Field: "telegram.parse_mode",
			Msg:   fmt.Sprintf("parse_mode must be HTML or Markdown, got %q", cfg.ParseMode),
		}
	}
	return nil
}
