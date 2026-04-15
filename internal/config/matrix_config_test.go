package config

import (
	"os"
	"testing"
)

func writeMatrixConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_MatrixDefaults(t *testing.T) {
	path := writeMatrixConfig(t, "")
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Matrix.Enabled {
		t.Error("expected matrix disabled by default")
	}
	if cfg.Matrix.MsgType != "m.text" {
		t.Errorf("expected msg_type m.text, got %s", cfg.Matrix.MsgType)
	}
}

func TestLoad_MatrixSection(t *testing.T) {
	path := writeMatrixConfig(t, `
[matrix]
enabled = true
homeserver = "https://matrix.example.com"
access_token = "syt_abc123"
room_id = "!roomid:example.com"
msg_type = "m.notice"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Matrix.Enabled {
		t.Error("expected matrix enabled")
	}
	if cfg.Matrix.Homeserver != "https://matrix.example.com" {
		t.Errorf("unexpected homeserver: %s", cfg.Matrix.Homeserver)
	}
	if cfg.Matrix.MsgType != "m.notice" {
		t.Errorf("unexpected msg_type: %s", cfg.Matrix.MsgType)
	}
}

func TestLoad_MatrixMissingHomeserver(t *testing.T) {
	path := writeMatrixConfig(t, `
[matrix]
enabled = true
access_token = "tok"
room_id = "!room:example.com"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing homeserver")
	}
}

func TestLoad_MatrixInvalidMsgType(t *testing.T) {
	path := writeMatrixConfig(t, `
[matrix]
enabled = true
homeserver = "https://matrix.example.com"
access_token = "tok"
room_id = "!room:example.com"
msg_type = "m.bad"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid msg_type")
	}
}
