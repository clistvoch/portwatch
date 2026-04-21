package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeRevoltConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return f.Name()
}

func TestLoad_RevoltDefaults(t *testing.T) {
	path := writeRevoltConfig(t, `
[revolt]
enabled = false
`)
	cfg, err := Load(path)
	require.NoError(t, err)
	assert.False(t, cfg.Revolt.Enabled)
	assert.Equal(t, "portwatch alert", cfg.Revolt.Username)
	assert.Equal(t, "", cfg.Revolt.WebhookURL)
}

func TestLoad_RevoltSection(t *testing.T) {
	path := writeRevoltConfig(t, `
[revolt]
enabled = true
webhook_url = "https://revolt.example.com/hook/abc123"
username = "portbot"
`)
	cfg, err := Load(path)
	require.NoError(t, err)
	assert.True(t, cfg.Revolt.Enabled)
	assert.Equal(t, "https://revolt.example.com/hook/abc123", cfg.Revolt.WebhookURL)
	assert.Equal(t, "portbot", cfg.Revolt.Username)
}

func TestLoad_RevoltMissingWebhook(t *testing.T) {
	path := writeRevoltConfig(t, `
[revolt]
enabled = true
`)
	_, err := Load(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "revolt")
}
