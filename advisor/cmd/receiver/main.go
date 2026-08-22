package main

import (
	"log"
	"net/http"

	"github.com/lucasschilin/s5n/advisor/internal/adapter/queue"
	"github.com/lucasschilin/s5n/advisor/internal/adapter/whatsapp"
	"github.com/lucasschilin/s5n/advisor/internal/config"
)

func init() {
	config.Load()
}
func main() {
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is healthy"))
	})

	// Initialize Queue producer
	queueProducer, err := queue.NewRabbitMQProducer(
		config.AppConfig.RabbitMQConnURL,
		config.AppConfig.IncomingMessagesQueueName,
	)
	if err != nil {
		log.Fatalf("Error initializing queue producer: %v", err)
	}
	defer queueProducer.Close()

	// Initialize WhatsappWebhook handler
	webhookHandler := whatsapp.NewWebhookHandler(queueProducer)

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			webhookHandler.HandleWebhook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})

	port := config.AppConfig.AppPort
	log.Printf("🚀 Server running on port :%s", port)
	log.Printf(
		"📌 GET route for Webhook ready at: http://localhost:%s/webhook", port,
	)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
