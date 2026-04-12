package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildLokiChange(port int, t monitor.ChangeType) monitor.Change {
	return monitor.Change{Port: port, Proto: "tcp", Type: t}
}

func TestLokiHandler_NoChanges(t *testing.T) {
	h := NewLokiHandler("http://localhost:3100", nil)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestLokiHandler_SendsPayload(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/loki/api/v1/push" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var err error
		received, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	h := NewLokiHandler(ts.URL, map[string]string{"job": "portwatch", "host": "test"})
	changes := []monitor.Change{
		buildLokiChange(8080, monitor.Opened),
		buildLokiChange(9090, monitor.Closed),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload lokiPayload
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(payload.Streams) != 1 {
		t.Fatalf("expected 1 stream, got %d", len(payload.Streams))
	}
	if len(payload.Streams[0].Values) != 2 {
		t.Errorf("expected 2 log values, got %d", len(payload.Streams[0].Values))
	}
}

func TestLokiHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewLokiHandler(ts.URL, nil)
	err := h.Handle([]monitor.Change{buildLokiChange(80, monitor.Opened)})
	if err == nil {
		t.Fatal("expected error on server 500, got nil")
	}
}

func TestLokiHandler_DefaultLabels(t *testing.T) {
	h := NewLokiHandler("http://localhost:3100", nil)
	if h.labels["job"] != "portwatch" {
		t.Errorf("expected default job label 'portwatch', got %q", h.labels["job"])
	}
}
