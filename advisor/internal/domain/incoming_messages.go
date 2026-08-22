package domain

import (
	"context"
	"time"
)

// RawIncomingMessage represents the rawwebhook payload saved immediately to the queue.
type RawIncomingMessage struct {
	MessageID  string    `json:"id"`          // message id
	UserID     string    `json:"user_id"`     // user id (ex: BR.13491208655302741918)
	Body       string    `json:"body"`        // message body
	ReceivedAt time.Time `json:"received_at"` // message received time
}

// IncomingMessagesRecorder defines the contract for recording incoming messages.
type IncomingMessagesQueueProducer interface {
	EnqueueIncomingMessage(
		ctx context.Context, message *RawIncomingMessage,
	) error
}
