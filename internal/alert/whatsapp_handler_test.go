package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildWhatsAppChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Port: monitor.PortInfo{Port: port, Proto: "tcp"},
		Kind: kind,
	}
}

func TestWhatsAppHandler_NoChanges(t *testing.T) {
	h := NewWhatsAppHandler("tok", "123", "447000000000", "https://example.com")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWhatsAppHandler_SendsPayload(t *testing.T) {
	var gotAuth, gotBody string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"messages":[{"id":"wamid.1"}]}`))
	}))
	defer srv.Close()

	h := NewWhatsAppHandler("mytoken", "987654321", "447000000000", srv.URL)
	changes := []monitor.Change{
		buildWhatsAppChange(8080, monitor.Opened),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(gotAuth, "Bearer mytoken") {
		t.Errorf("expected Bearer token in Authorization, got %q", gotAuth)
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(gotBody), &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if payload["messaging_product"] != "whatsapp" {
		t.Errorf("expected messaging_product=whatsapp, got %v", payload["messaging_product"])
	}
	textBlock, ok := payload["text"].(map[string]any)
	if !ok {
		t.Fatal("missing text block in payload")
	}
	if body, _ := textBlock["body"].(string); !strings.Contains(body, "8080") {
		t.Errorf("expected port 8080 in message body, got %q", body)
	}
}

func TestWhatsAppHandler_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	h := NewWhatsAppHandler("tok", "123", "447000000000", srv.URL)
	err := h.Handle([]monitor.Change{buildWhatsAppChange(443, monitor.Closed)})
	if err == nil {
		t.Fatal("expected error on 500 response")
	}
}
