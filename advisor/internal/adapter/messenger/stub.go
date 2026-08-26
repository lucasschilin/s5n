package messenger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lucasschilin/s5n/advisor/internal/domain"
)

type StubMessengerAdapter struct {
	queue domain.IncomingMessagesQueueProducer
}

func NewStubMessengerAdapter() *StubMessengerAdapter {
	return &StubMessengerAdapter{}
}

func (h *StubMessengerAdapter) HandleWebhook(queue domain.IncomingMessagesQueueProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			MessageID string `json:"message_id"`
			From      string `json:"from"`
			Body      string `json:"body"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		rawMessage := &domain.RawIncomingMessage{
			MessageID:  payload.MessageID,
			UserID:     payload.From,
			Body:       payload.Body,
			ReceivedAt: time.Now().UTC(),
		}

		log.Printf("Received message: %+v", rawMessage)

		if err := queue.EnqueueIncomingMessage(r.Context(), rawMessage); err != nil {
			http.Error(w, "failed to enqueue message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *StubMessengerAdapter) SendMessage(userID, body string, replyMessageID string) error {
	if replyMessageID == "" {
		fmt.Printf("Sending message to user [%s]: '%s'\n\n", userID, body)
		return nil
	}

	fmt.Printf("Answering message [%s] to user [%s]: '%s'\n\n", replyMessageID, userID, body)
	return nil
}
