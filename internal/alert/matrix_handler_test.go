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

func buildMatrixChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{Port: port, Kind: kind}
}

func TestMatrixHandler_NoChanges(t *testing.T) {
	h := NewMatrixHandler("https://matrix.example.com", "tok", "!room:example.com", "m.text")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMatrixHandler_SendsPayload(t *testing.T) {
	var gotBody map[string]string
	var gotAuth string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		raw, _ := io.ReadAll(r.Body)
		json.Unmarshal(raw, &gotBody)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"event_id":"$abc"}`))
	}))
	defer ts.Close()

	h := NewMatrixHandler(ts.URL, "mytoken", "!room:example.com", "m.notice")
	changes := []monitor.Change{
		buildMatrixChange(8080, monitor.Opened),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "Bearer mytoken" {
		t.Errorf("expected Bearer mytoken, got %s", gotAuth)
	}
	if gotBody["msgtype"] != "m.notice" {
		t.Errorf("expected msgtype m.notice, got %s", gotBody["msgtype"])
	}
	if !strings.Contains(gotBody["body"], "8080") {
		t.Errorf("expected body to mention port 8080, got: %s", gotBody["body"])
	}
}

func TestMatrixHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	h := NewMatrixHandler(ts.URL, "badtoken", "!room:example.com", "m.text")
	err := h.Handle([]monitor.Change{buildMatrixChange(22, monitor.Opened)})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}
