package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wperron/portwatch/internal/monitor"
)

func buildSPChange(t monitor.ChangeType) monitor.Change {
	return monitor.Change{
		Type: t,
		Port: monitor.PortInfo{Port: 8080, Proto: "tcp"},
	}
}

func TestStatusPageHandler_NoChanges(t *testing.T) {
	h := NewStatusPageHandler("key", "page", "comp", "http://localhost")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestStatusPageHandler_SendsPayload(t *testing.T) {
	var gotBody map[string]any
	var gotAuth string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := NewStatusPageHandler("mykey", "pageid", "compid", srv.URL)
	if err := h.Handle([]monitor.Change{buildSPChange(monitor.Opened)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotAuth != "OAuth mykey" {
		t.Errorf("expected OAuth mykey, got %q", gotAuth)
	}

	comp, ok := gotBody["component"].(map[string]any)
	if !ok {
		t.Fatal("missing component key in payload")
	}
	if comp["status"] != "under_maintenance" {
		t.Errorf("expected under_maintenance, got %v", comp["status"])
	}
}

func TestStatusPageHandler_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	h := NewStatusPageHandler("k", "p", "c", srv.URL)
	err := h.Handle([]monitor.Change{buildSPChange(monitor.Closed)})
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestStatusPageHandler_InvalidURL(t *testing.T) {
	h := NewStatusPageHandler("k", "p", "c", "://bad-url")
	err := h.Handle([]monitor.Change{buildSPChange(monitor.Opened)})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
