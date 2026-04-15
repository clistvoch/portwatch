package config

import "fmt"

type MatrixConfig struct {
	Enabled    bool   `toml:"enabled"`
	Homeserver string `toml:"homeserver"`
	Token      string `toml:"access_token"`
	RoomID     string `toml:"room_id"`
	MsgType    string `toml:"msg_type"`
}

func defaultMatrixConfig() MatrixConfig {
	return MatrixConfig{
		Enabled: false,
		MsgType: "m.text",
	}
}

func validateMatrix(c MatrixConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Homeserver == "" {
		return &ValidationError{Field: "matrix.homeserver", Msg: "homeserver is required"}
	}
	if c.Token == "" {
		return &ValidationError{Field: "matrix.access_token", Msg: "access_token is required"}
	}
	if c.RoomID == "" {
		return &ValidationError{Field: "matrix.room_id", Msg: "room_id is required"}
	}
	allowed := map[string]bool{"m.text": true, "m.notice": true}
	if !allowed[c.MsgType] {
		return &ValidationError{Field: "matrix.msg_type", Msg: fmt.Sprintf("invalid msg_type %q: must be m.text or m.notice", c.MsgType)}
	}
	return nil
}
