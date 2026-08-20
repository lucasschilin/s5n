package main

import (
	"log"
	"net/http"

	"github.com/lucasschilin/s5n/advisor/internal/config"
)

func init() {
	config.Load()
}
func main() {
	port := config.AppConfig.AppPort

	// Rota básica de saúde
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is healthy"))
	})

	log.Printf("🚀 Server running on port :%s", port)
	log.Printf("📌 GET route for Webhook ready at: http://localhost:%s/webhook", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
