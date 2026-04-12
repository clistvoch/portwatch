package alert

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

// fakeKafkaProducer captures sent messages for assertions.
type fakeKafkaProducer struct {
	messages [][]byte
	errSend  error
}

func (f *fakeKafkaProducer) SendMessage(_ string, payload []byte) error {
	if f.errSend != nil {
		return f.errSend
	}
	f.messages = append(f.messages, payload)
	return nil
}

func (f *fakeKafkaProducer) Close() error { return nil }

func buildKafkaChange(port int, t monitor.ChangeType) monitor.Change {
	return monitor.Change{Port: port, Type: t}
}

func TestKafkaHandler_NoChanges(t *testing.T) {
	prod := &fakeKafkaProducer{}
	h := NewKafkaHandler(prod, "portwatch-alerts")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prod.messages) != 0 {
		t.Errorf("expected no messages, got %d", len(prod.messages))
	}
}

func TestKafkaHandler_SendsPayload(t *testing.T) {
	prod := &fakeKafkaProducer{}
	h := NewKafkaHandler(prod, "portwatch-alerts")

	changes := []monitor.Change{
		buildKafkaChange(8080, monitor.Opened),
		buildKafkaChange(22, monitor.Closed),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prod.messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(prod.messages))
	}

	var p kafkaPayload
	if err := json.Unmarshal(prod.messages[0], &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(p.Changes) != 2 {
		t.Errorf("expected 2 changes in payload, got %d", len(p.Changes))
	}
	if p.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestKafkaHandler_ProducerError(t *testing.T) {
	prod := &fakeKafkaProducer{errSend: errors.New("broker unavailable")}
	h := NewKafkaHandler(prod, "portwatch-alerts")

	changes := []monitor.Change{buildKafkaChange(443, monitor.Opened)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error from producer, got nil")
	}
}
