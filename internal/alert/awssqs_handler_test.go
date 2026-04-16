package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildSQSChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{Port: port, Kind: kind}
}

func TestAWSSQSHandler_NoChanges(t *testing.T) {
	h := NewAWSSQSHandler("http://example.com", "us-east-1", "k", "s", "")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAWSSQSHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewAWSSQSHandler(ts.URL, "us-east-1", "key", "secret", "portwatch")
	changes := []monitor.Change{buildSQSChange(8080, monitor.Opened)}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["source"] != "portwatch" {
		t.Errorf("unexpected source: %v", received["source"])
	}
}

func TestAWSSQSHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewAWSSQSHandler(ts.URL, "us-east-1", "key", "secret", "")
	err := h.Handle([]monitor.Change{buildSQSChange(9090, monitor.Closed)})
	if err == nil {
		t.Error("expected error on 500 response")
	}
}

func TestAWSSQSHandler_InvalidURL(t *testing.T) {
	h := NewAWSSQSHandler("http://127.0.0.1:0", "us-east-1", "k", "s", "")
	err := h.Handle([]monitor.Change{buildSQSChange(1234, monitor.Opened)})
	if err == nil {
		t.Error("expected connection error")
	}
}
