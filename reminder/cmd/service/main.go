package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucasschilin/s5n/reminder/internal/adapter/queue"
	"github.com/lucasschilin/s5n/reminder/internal/config"
	"github.com/lucasschilin/s5n/reminder/internal/domain"
)

func init() {
	config.Load()
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer, err := queue.NewRabbitMQConsumer[domain.ReminderQueuePayload](
		config.AppConfig.RabbitMQConnURL, "reminder_agent",
	)
	if err != nil {
		log.Fatalf("❌ Error initializing queue consumer: %v", err)
	}
	defer consumer.Close()

	printActionFunc := func(ctx context.Context, msg domain.ReminderQueuePayload) error {
		log.Printf("📩 Received message: %+v", msg)
		return nil
	}

	err = consumer.StartConsuming(ctx, printActionFunc)
	if err != nil {
		log.Fatalf("❌ Error starting consumption: %v", err)
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	log.Println("🟢 Intent Router Service running. Press CTRL+C to exit.")
	<-stopSignal

	log.Println("⏳ Finishing up, waiting for ongoing message processing to complete...")
	cancel()
	time.Sleep(1 * time.Second)
	log.Println("👋 Intent Router Service finished successfully.")
}
