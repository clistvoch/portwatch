package alert

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func buildHMACChange(kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Port: monitor.PortInfo{Port: 9090, Proto: "tcp"},
	}
}

func TestWebhookHMACHandler_NoChanges(t *testing.T) {
	h := NewWebhookHMACHandler("http://example.com", "secret", "sha256", "X-Sig")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error for empty changes, got %v", err)
	}
}

func TestWebhookHMACHandler_SendsPayload(t *testing.T) {
	var gotSig, gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig = r.Header.Get("X-Portwatch-Signature")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	secret := "topsecret"
	h := NewWebhookHMACHandler(ts.URL, secret, "sha256", "X-Portwatch-Signature")
	changes := []monitor.Change{buildHMACChange(monitor.Opened)}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// verify signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(gotBody))
	expected := hex.EncodeToString(mac.Sum(nil))
	if gotSig != expected {
		t.Errorf("signature mismatch: got %s, want %s", gotSig, expected)
	}
}

func TestWebhookHMACHandler_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewWebhookHMACHandler(ts.URL, "secret", "sha256", "X-Sig")
	err := h.Handle([]monitor.Change{buildHMACChange(monitor.Opened)})
	if err == nil {
		t.Error("expected error on 500 response")
	}
}

func TestWebhookHMACHandler_InvalidURL(t *testing.T) {
	h := NewWebhookHMACHandler("://bad-url", "secret", "sha256", "X-Sig")
	err := h.Handle([]monitor.Change{buildHMACChange(monitor.Opened)})
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}
