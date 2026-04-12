package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildESChange(kind monitor.ChangeKind, port int) monitor.Change {
	return monitor.Change{Kind: kind, Port: port, Proto: "tcp"}
}

func TestElasticsearchHandler_NoChanges(t *testing.T) {
	h := NewElasticsearchHandler("http://localhost:9200", "portwatch", "", "", 5)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestElasticsearchHandler_SendsPayload(t *testing.T) {
	var received []map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var doc map[string]interface{}
		_ = json.Unmarshal(body, &doc)
		received = append(received, doc)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	h := NewElasticsearchHandler(srv.URL, "portwatch", "", "", 5)
	changes := []monitor.Change{
		buildESChange(monitor.Opened, 8080),
		buildESChange(monitor.Closed, 22),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(received))
	}
	if received[0]["port"].(float64) != 8080 {
		t.Errorf("unexpected port in first doc: %v", received[0]["port"])
	}
	if received[1]["kind"] != "closed" {
		t.Errorf("unexpected kind in second doc: %v", received[1]["kind"])
	}
}

func TestElasticsearchHandler_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	h := NewElasticsearchHandler(srv.URL, "portwatch", "", "", 5)
	err := h.Handle([]monitor.Change{buildESChange(monitor.Opened, 443)})
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestElasticsearchHandler_InvalidURL(t *testing.T) {
	h := NewElasticsearchHandler("http://127.0.0.1:1", "portwatch", "", "", 1)
	err := h.Handle([]monitor.Change{buildESChange(monitor.Opened, 80)})
	if err == nil {
		t.Fatal("expected connection error")
	}
}
