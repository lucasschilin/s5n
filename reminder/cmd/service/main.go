package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucasschilin/s5n/reminder/internal/adapter/database"
	"github.com/lucasschilin/s5n/reminder/internal/adapter/queue"
	"github.com/lucasschilin/s5n/reminder/internal/config"
	"github.com/lucasschilin/s5n/reminder/internal/domain"
	"github.com/lucasschilin/s5n/reminder/internal/services"
)

func init() {
	config.Load()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connection (no automatic migrations)
	db, err := database.NewSQLiteDB("reminders.db")
	if err != nil {
		log.Fatalf("❌ Error initializing database: %v", err)
	}
	defer db.Close()

	// 2. Wire Repository -> Service
	reminderRepo := database.NewSQLiteRepository(db)
	reminderService := services.NewReminderService(reminderRepo)

	// 3. Initialize RabbitMQ Consumer
	consumer, err := queue.NewRabbitMQConsumer[domain.ReminderQueuePayload](
		config.AppConfig.RabbitMQConnURL, "reminder_agent",
	)
	if err != nil {
		log.Fatalf("❌ Error initializing queue consumer: %v", err)
	}
	defer consumer.Close()

	// 4. Start processing messages through domain service
	err = consumer.StartConsuming(ctx, reminderService.ProcessQueueMessage)
	if err != nil {
		log.Fatalf("❌ Error starting queue consumer: %v", err)
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	log.Println("🟢 Service running. Press CTRL+C to exit.")
	<-stopSignal

	log.Println("⏳ Shutting down gracefully...")
	cancel()
	time.Sleep(1 * time.Second)
	log.Println("👋 Service stopped successfully.")
}
