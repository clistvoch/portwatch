package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildGPSChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{Port: port, Kind: kind}
}

func TestGooglePubSubHandler_NoChanges(t *testing.T) {
	h := NewGooglePubSubHandler("proj", "topic")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestGooglePubSubHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode error: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewGooglePubSubHandler("myproject", "mytopic")
	h.baseURL = ts.URL

	changes := []monitor.Change{
		buildGPSChange(8080, monitor.Opened),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received == nil {
		t.Fatal("expected payload to be received")
	}
}

func TestGooglePubSubHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewGooglePubSubHandler("proj", "topic")
	h.baseURL = ts.URL

	changes := []monitor.Change{buildGPSChange(443, monitor.Closed)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error on 500 response")
	}
}
