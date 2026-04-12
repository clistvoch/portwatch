package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeDiscordConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return f.Name()
}

func TestLoad_DiscordDefaults(t *testing.T) {
	path := writeDiscordConfig(t, `
[scan]
range = "1-1024"
`)
	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "", cfg.Discord.WebhookURL)
	assert.Equal(t, "portwatch", cfg.Discord.Username)
	assert.Equal(t, true, cfg.Discord.Enabled)
}

func TestLoad_DiscordSection(t *testing.T) {
	path := writeDiscordConfig(t, `
[scan]
range = "1-1024"

[discord]
webhook_url = "https://discord.com/api/webhooks/123/abc"
username = "alertbot"
enabled = true
`)
	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "https://discord.com/api/webhooks/123/abc", cfg.Discord.WebhookURL)
	assert.Equal(t, "alertbot", cfg.Discord.Username)
	assert.True(t, cfg.Discord.Enabled)
}

func TestLoad_DiscordMissingWebhook(t *testing.T) {
	path := writeDiscordConfig(t, `
[scan]
range = "1-1024"

[discord]
enabled = true
`)
	cfg, err := Load(path)
	require.NoError(t, err)
	err = Validate(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "discord")
}
