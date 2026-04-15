package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wiring/portwatch/internal/monitor"
)

func buildJiraChange(kind monitor.ChangeKind, port int) monitor.Change {
	return monitor.Change{Kind: kind, Port: monitor.PortInfo{Port: port, Proto: "tcp"}}
}

func TestJiraHandler_NoChanges(t *testing.T) {
	h := NewJiraHandler("http://localhost", "u", "tok", "OPS", "Bug", "High")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error for empty changes, got %v", err)
	}
}

func TestJiraHandler_SendsPayload(t *testing.T) {
	var received map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		user, pass, ok := r.BasicAuth()
		if !ok || user != "user" || pass != "secret" {
			t.Errorf("bad basic auth: user=%s pass=%s", user, pass)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	h := NewJiraHandler(ts.URL, "user", "secret", "OPS", "Bug", "High")
	changes := []monitor.Change{
		buildJiraChange(monitor.Opened, 8080),
		buildJiraChange(monitor.Closed, 22),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fields, ok := received["fields"].(map[string]any)
	if !ok {
		t.Fatal("missing fields in payload")
	}
	summary, _ := fields["summary"].(string)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestJiraHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewJiraHandler(ts.URL, "u", "tok", "OPS", "Bug", "High")
	changes := []monitor.Change{buildJiraChange(monitor.Opened, 443)}
	if err := h.Handle(changes); err == nil {
		t.Error("expected error on server 500")
	}
}

func TestJiraHandler_InvalidURL(t *testing.T) {
	h := NewJiraHandler("http://127.0.0.1:1", "u", "tok", "OPS", "Bug", "High")
	changes := []monitor.Change{buildJiraChange(monitor.Opened, 80)}
	if err := h.Handle(changes); err == nil {
		t.Error("expected error for unreachable server")
	}
}
