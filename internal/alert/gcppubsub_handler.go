package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"

	"github.com/patrickward/portwatch/internal/monitor"
)

// GCPPubSubHandler publishes port change events to a GCP Pub/Sub topic.
type GCPPubSubHandler struct {
	client  *pubsub.Client
	topic   *pubsub.Topic
	project string
	topicID string
}

type pubSubPayload struct {
	Timestamp string          `json:"timestamp"`
	Changes   []monitor.Change `json:"changes"`
}

// NewGCPPubSubHandler creates a new GCPPubSubHandler.
func NewGCPPubSubHandler(projectID, topicID, credFile string) (*GCPPubSubHandler, error) {
	ctx := context.Background()
	var opts []option.ClientOption
	if credFile != "" {
		opts = append(opts, option.WithCredentialsFile(credFile))
	}
	client, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("gcppubsub: create client: %w", err)
	}
	return &GCPPubSubHandler{
		client:  client,
		topic:   client.Topic(topicID),
		project: projectID,
		topicID: topicID,
	}, nil
}

// Handle sends changes to GCP Pub/Sub.
func (h *GCPPubSubHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	payload := pubSubPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Changes:   changes,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gcppubsub: marshal: %w", err)
	}
	ctx := context.Background()
	result := h.topic.Publish(ctx, &pubsub.Message{Data: data})
	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("gcppubsub: publish: %w", err)
	}
	return nil
}

// Close stops the topic and closes the client.
func (h *GCPPubSubHandler) Close() error {
	h.topic.Stop()
	return h.client.Close()
}
