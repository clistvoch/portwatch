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

func buildRevoltChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: port,
		Kind: kind,
		Proto: "tcp",
	}
}

func TestRevoltHandler_NoChanges(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := alert.NewRevoltHandler(ts.URL, "portbot")
	err := h.Handle([]monitor.Change{})
	require.NoError(t, err)
	assert.False(t, called, "should not send request when no changes")
}

func TestRevoltHandler_SendsPayload(t *testing.T) {
	var body map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := alert.NewRevoltHandler(ts.URL, "portbot")
	changes := []monitor.Change{
		buildRevoltChange(8080, monitor.ChangeOpened),
	}
	err := h.Handle(changes)
	require.NoError(t, err)
	assert.Equal(t, "portbot", body["username"])
	content, ok := body["content"].(string)
	require.True(t, ok)
	assert.Contains(t, content, "8080")
}

func TestRevoltHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := alert.NewRevoltHandler(ts.URL, "portbot")
	changes := []monitor.Change{
		buildRevoltChange(9090, monitor.ChangeClosed),
	}
	err := h.Handle(changes)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}
