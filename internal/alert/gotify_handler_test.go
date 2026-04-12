package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildGotifyChange(kind monitor.ChangeKind, port uint16) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
	}
}

func TestGotifyHandler_NoChanges(t *testing.T) {
	h := NewGotifyHandler("http://localhost", "tok", 5)
	if err := h.Handle(nil); err != nil {
		t.Errorf("expected no error on empty changes, got %v", err)
	}
}

func TestGotifyHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewGotifyHandler(ts.URL, "mytoken", 7)
	changes := []monitor.Change{
		buildGotifyChange(monitor.Opened, 8080),
		buildGotifyChange(monitor.Closed, 9090),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["title"] == nil {
		t.Error("expected title in payload")
	}
	if received["priority"].(float64) != 7 {
		t.Errorf("expected priority 7, got %v", received["priority"])
	}
	if received["message"] == nil {
		t.Error("expected message in payload")
	}
}

func TestGotifyHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewGotifyHandler(ts.URL, "tok", 5)
	err := h.Handle([]monitor.Change{buildGotifyChange(monitor.Opened, 443)})
	if err == nil {
		t.Fatal("expected error on server 500")
	}
}
