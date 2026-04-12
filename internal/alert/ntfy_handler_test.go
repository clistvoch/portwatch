package alert

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildNtfyChange(port uint16, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: port,
		Kind: kind,
	}
}

func TestNtfyHandler_NoChanges(t *testing.T) {
	h := NewNtfyHandler("", "alerts", 3)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error for empty changes, got %v", err)
	}
}

func TestNtfyHandler_SendsPayload(t *testing.T) {
	var gotBody, gotTitle, gotPriority string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTitle = r.Header.Get("Title")
		gotPriority = r.Header.Get("Priority")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewNtfyHandler(ts.URL, "portwatch", 4)
	changes := []monitor.Change{
		buildNtfyChange(8080, monitor.Opened),
		buildNtfyChange(9090, monitor.Closed),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(gotTitle, "2 port change") {
		t.Errorf("expected title to mention 2 changes, got %q", gotTitle)
	}
	if gotPriority != "4" {
		t.Errorf("expected priority 4, got %q", gotPriority)
	}
	if !strings.Contains(gotBody, "8080") || !strings.Contains(gotBody, "9090") {
		t.Errorf("expected body to contain port numbers, got %q", gotBody)
	}
}

func TestNtfyHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewNtfyHandler(ts.URL, "portwatch", 3)
	err := h.Handle([]monitor.Change{buildNtfyChange(443, monitor.Opened)})
	if err == nil {
		t.Fatal("expected error on server 500, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to mention status 500, got %v", err)
	}
}

func TestNtfyHandler_DefaultPriority(t *testing.T) {
	h := NewNtfyHandler("", "alerts", 0)
	if h.priority != 3 {
		t.Errorf("expected default priority 3, got %d", h.priority)
	}
}
