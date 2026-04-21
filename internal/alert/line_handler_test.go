package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/monitor"
)

func buildLineChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestLineHandler_NoChanges(t *testing.T) {
	h := alert.NewLineHandler("tok", "https://example.com", "prefix")
	err := h.Handle([]monitor.Change{})
	assert.NoError(t, err)
}

func TestLineHandler_SendsPayload(t *testing.T) {
	var gotForm string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer tok", r.Header.Get("Authorization"))
		require.NoError(t, r.ParseForm())
		gotForm = r.FormValue("message")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": 200, "message": "ok"})
	}))
	defer ts.Close()

	h := alert.NewLineHandler("tok", ts.URL, "[portwatch]")
	err := h.Handle([]monitor.Change{
		buildLineChange(8080, monitor.ChangeOpened),
	})
	require.NoError(t, err)
	assert.Contains(t, gotForm, "8080")
	assert.Contains(t, gotForm, "[portwatch]")
}

func TestLineHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := alert.NewLineHandler("tok", ts.URL, "prefix")
	err := h.Handle([]monitor.Change{
		buildLineChange(9090, monitor.ChangeClosed),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}
