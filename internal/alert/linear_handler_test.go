package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/username/portwatch/internal/monitor"
)

func buildLinearChange(kind monitor.ChangeKind, port int) monitor.Change {
	return monitor.Change{Kind: kind, Port: monitor.PortInfo{Port: port, Proto: "tcp"}}
}

func TestLinearHandler_NoChanges(t *testing.T) {
	h := NewLinearHandler("key", "team1", "", "https://api.linear.app", "", 2, nil)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestLinearHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "test-key" {
			t.Errorf("missing or wrong Authorization header")
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	h := NewLinearHandler("test-key", "team-abc", "proj-1", ts.URL, "", 1, []string{"label-x"})
	changes := []monitor.Change{
		buildLinearChange(monitor.Opened, 8080),
		buildLinearChange(monitor.Closed, 22),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["teamId"] != "team-abc" {
		t.Errorf("expected teamId team-abc, got %v", received["teamId"])
	}
	if received["projectId"] != "proj-1" {
		t.Errorf("expected projectId proj-1, got %v", received["projectId"])
	}
	title, _ := received["title"].(string)
	if title == "" {
		t.Error("expected non-empty title")
	}
}

func TestLinearHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h", "team1", "", ts.URL, "", 2, nil)
	err := h.Handle([]monitor.Change{buildLinearChange(monitor.Opened, 9090)})
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestLinearHandler_InvalidURL(t *testing.T) {
	h := NewLinearHandler("key", "team1", "", "://bad-url", "", 2, nil)
	err := h.Handle([]monitor.Change{buildLinearChange(monitor.Opened, 443)})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
