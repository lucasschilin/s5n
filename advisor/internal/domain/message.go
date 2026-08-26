package domain

import (
	"context"
	"net/http"
	"time"
)

// RawIncomingMessage represents the rawwebhook payload saved immediately to the queue.
type RawIncomingMessage struct {
	MessageID  string    `json:"message_id"`  // message id
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

type RawOutgoingMessage struct {
	ReplyMessageID string `json:"reply_message_id"` // reply message id
	UserID         string `json:"user_id"`          // user id (ex: BR.13491208655302741918)
	Body           string `json:"body"`             // message body
}

type Messenger interface {
	HandleWebhook(queue IncomingMessagesQueueProducer) http.HandlerFunc
	SendMessage(userID, body string, replyMessageID string) error
}
