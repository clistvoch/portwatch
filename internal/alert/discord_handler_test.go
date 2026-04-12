package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func buildDiscordChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: port,
		Kind: kind,
		Proto: "tcp",
	}
}

func TestDiscordHandler_NoChanges(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := alert.NewDiscordHandler(ts.URL, "portwatch")
	err := h.Handle([]monitor.Change{})
	require.NoError(t, err)
	assert.False(t, called, "should not POST when there are no changes")
}

func TestDiscordHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&received))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	h := alert.NewDiscordHandler(ts.URL, "portwatch")
	changes := []monitor.Change{
		buildDiscordChange(8080, monitor.Opened),
		buildDiscordChange(9090, monitor.Closed),
	}
	err := h.Handle(changes)
	require.NoError(t, err)

	assert.Equal(t, "portwatch", received["username"])
	content, ok := received["content"].(string)
	require.True(t, ok)
	assert.Contains(t, content, "8080")
	assert.Contains(t, content, "9090")
}

func TestDiscordHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := alert.NewDiscordHandler(ts.URL, "portwatch")
	err := h.Handle([]monitor.Change{buildDiscordChange(443, monitor.Opened)})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}
