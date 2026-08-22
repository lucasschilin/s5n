package whatsapp

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/lucasschilin/s5n/advisor/internal/domain"
)

type WebhookHandler struct {
	queue domain.QueueProducer
}

func NewWebhookHandler(queue domain.QueueProducer) *WebhookHandler {
	return &WebhookHandler{
		queue: queue,
	}
}

func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		MessageID string `json:"id"`
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

	if err := h.queue.EnqueueIncomingMessage(r.Context(), rawMessage); err != nil {
		http.Error(w, "failed to enqueue message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
