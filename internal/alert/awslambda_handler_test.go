package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildLambdaChange(port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{Port: port, Kind: kind}
}

func TestAWSLambdaHandler_NoChanges(t *testing.T) {
	h := NewAWSLambdaHandler("my-fn", "us-east-1", "key", "secret", "Event", 5)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAWSLambdaHandler_SendsPayload(t *testing.T) {
	var received lambdaPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	h := NewAWSLambdaHandler("fn", "us-east-1", "key", "secret", "Event", 5).(*awsLambdaHandler)
	h.endpointURL = ts.URL

	changes := []monitor.Change{buildLambdaChange(8080, monitor.Opened)}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(received.Changes))
	}
	if received.Event != "portwatch.change" {
		t.Errorf("expected event portwatch.change, got %s", received.Event)
	}
}

func TestAWSLambdaHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewAWSLambdaHandler("fn", "us-east-1", "key", "secret", "Event", 5).(*awsLambdaHandler)
	h.endpointURL = ts.URL

	if err := h.Handle([]monitor.Change{buildLambdaChange(9090, monitor.Closed)}); err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestAWSLambdaHandler_InvalidURL(t *testing.T) {
	h := NewAWSLambdaHandler("fn", "us-east-1", "key", "secret", "Event", 5).(*awsLambdaHandler)
	h.endpointURL = "http://127.0.0.1:1"

	if err := h.Handle([]monitor.Change{buildLambdaChange(1234, monitor.Opened)}); err == nil {
		t.Fatal("expected connection error")
	}
}
