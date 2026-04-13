package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func buildRocketChatChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{Port: port, Kind: kind}
}

func TestRocketChatHandler_NoChanges(t *testing.T) {
	h := alert.NewRocketChatHandler("http://example.com/hook", nil)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected nil error for empty changes, got %v", err)
	}
}

func TestRocketChatHandler_SendsPayload(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	changes := []monitor.Change{
		buildRocketChatChange(8080, monitor.Opened),
		buildRocketChatChange(9090, monitor.Closed),
	}
	h := alert.NewRocketChatHandler(ts.URL, ts.Client())
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["text"] == "" {
		t.Error("expected non-empty text field in payload")
	}
}

func TestRocketChatHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	changes := []monitor.Change{buildRocketChatChange(443, monitor.Opened)}
	h := alert.NewRocketChatHandler(ts.URL, ts.Client())
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error on 500 response, got nil")
	}
}
