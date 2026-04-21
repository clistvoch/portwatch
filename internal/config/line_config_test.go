package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeLineConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return f.Name()
}

func TestLoad_LineDefaults(t *testing.T) {
	path := writeLineConfig(t, "[line]\nenabled = false\n")
	cfg, err := Load(path)
	require.NoError(t, err)
	assert.False(t, cfg.Line.Enabled)
	assert.Equal(t, "https://notify-api.line.me/api/notify", cfg.Line.APIURL)
	assert.Equal(t, "portwatch alert", cfg.Line.StickerPrefix)
}

func TestLoad_LineSection(t *testing.T) {
	path := writeLineConfig(t, `
[line]
enabled = true
token = "mytoken"
api_url = "https://custom.line.example.com/notify"
sticker_prefix = "[ALERT]"
`)
	cfg, err := Load(path)
	require.NoError(t, err)
	assert.True(t, cfg.Line.Enabled)
	assert.Equal(t, "mytoken", cfg.Line.Token)
	assert.Equal(t, "https://custom.line.example.com/notify", cfg.Line.APIURL)
	assert.Equal(t, "[ALERT]", cfg.Line.StickerPrefix)
}

func TestLoad_LineMissingToken(t *testing.T) {
	path := writeLineConfig(t, "[line]\nenabled = true\n")
	_, err := Load(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "line.token")
}
