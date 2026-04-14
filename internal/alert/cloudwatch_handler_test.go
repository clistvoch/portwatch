package alert_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

func buildCloudWatchChange(port int, status monitor.ChangeStatus) monitor.Change {
	return monitor.Change{Port: port, Status: status}
}

func TestCloudWatchHandler_NoChanges(t *testing.T) {
	h := alert.NewCloudWatchHandler("us-east-1", "PortWatch", "PortChange", "key", "secret")
	if err := h.Handle(context.Background(), nil); err != nil {
		t.Fatalf("expected no error on empty changes, got %v", err)
	}
}

func TestCloudWatchHandler_SendsPayload(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if r.Header.Get("X-Amz-Access-Key") == "" {
			t.Error("expected access key header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	h := alert.NewCloudWatchHandler("us-east-1", "PortWatch", "PortChange", "TESTKEY", "TESTSECRET")
	h.(*alert.CloudWatchHandler) // type assertion skipped; use exported setter in real code
	// Use internal endpoint override via a test-friendly constructor variant:
	h2 := alert.NewCloudWatchHandlerWithEndpoint("us-east-1", "PortWatch", "PortChange", "TESTKEY", "TESTSECRET", server.URL)

	changes := []monitor.Change{
		buildCloudWatchChange(8080, monitor.Opened),
		buildCloudWatchChange(9090, monitor.Closed),
	}
	if err := h2.Handle(context.Background(), changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["namespace"] != "PortWatch" {
		t.Errorf("unexpected namespace: %v", received["namespace"])
	}
	if received["value"].(float64) != 2 {
		t.Errorf("expected value 2, got %v", received["value"])
	}
}

func TestCloudWatchHandler_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	h := alert.NewCloudWatchHandlerWithEndpoint("us-east-1", "NS", "M", "k", "s", server.URL)
	changes := []monitor.Change{buildCloudWatchChange(22, monitor.Opened)}
	if err := h.Handle(context.Background(), changes); err == nil {
		t.Fatal("expected error on 500 response")
	}
}

func TestCloudWatchHandler_InvalidURL(t *testing.T) {
	h := alert.NewCloudWatchHandlerWithEndpoint("us-east-1", "NS", "M", "k", "s", "://bad-url")
	changes := []monitor.Change{buildCloudWatchChange(80, monitor.Opened)}
	if err := h.Handle(context.Background(), changes); err == nil {
		t.Fatal("expected error on invalid URL")
	}
}
