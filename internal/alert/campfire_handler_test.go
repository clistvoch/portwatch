package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmolesUC/portwatch/internal/monitor"
)

func buildCampfireChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{Port: port, Kind: kind}
}

func TestCampfireHandler_NoChanges(t *testing.T) {
	h := NewCampfireHandler("token", "123", "http://localhost")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCampfireHandler_SendsPayload(t *testing.T) {
	var gotBody map[string]interface{}
	var gotAuth string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	h := NewCampfireHandler("mytoken", "42", srv.URL)
	changes := []monitor.Change{
		buildCampfireChange(8080, monitor.Opened),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(gotAuth, "Basic") {
		t.Errorf("expected Basic auth header, got %q", gotAuth)
	}
	msg, ok := gotBody["message"].(map[string]interface{})
	if !ok {
		t.Fatal("expected message key in payload")
	}
	if msg["type"] != "TextMessage" {
		t.Errorf("expected TextMessage type, got %v", msg["type"])
	}
	if !strings.Contains(msg["body"].(string), "8080") {
		t.Errorf("expected port 8080 in body, got %v", msg["body"])
	}
}

func TestCampfireHandler_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	h := NewCampfireHandler("tok", "1", srv.URL)
	err := h.Handle([]monitor.Change{buildCampfireChange(9090, monitor.Closed)})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}
