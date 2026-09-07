package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucasschilin/s5n/advisor/internal/adapter/llm"
	"github.com/lucasschilin/s5n/advisor/internal/adapter/queue"
	"github.com/lucasschilin/s5n/advisor/internal/config"
	"github.com/lucasschilin/s5n/advisor/internal/domain"
	"github.com/lucasschilin/s5n/advisor/internal/service"
)

func init() {
	config.Load()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	classifierAdapter, err := llm.NewGeminiAdapter(
		context.Background(), config.AppConfig.GeminiAPIKey, config.AppConfig.GeminiModel,
	)
	if err != nil {
		log.Fatalf("❌ Error initializing LLM adapter: %v", err)
	}

	consumer, err := queue.NewRabbitMQConsumer[domain.RawIncomingMessage](
		config.AppConfig.RabbitMQConnURL, config.AppConfig.IncomingMessagesQueueName,
	)
	if err != nil {
		log.Fatalf("❌ Error initializing queue consumer: %v", err)
	}
	defer consumer.Close()

	producer, err := queue.NewRabbitMQProducer(config.AppConfig.RabbitMQConnURL, "")
	if err != nil {
		log.Fatalf("❌ Error initializing queue producer: %v", err)
	}
	defer producer.Close()

	routerService := service.NewIntentRouterService(classifierAdapter, producer, config.AppConfig.OutgoingMessagesQueueName)

	err = consumer.StartConsuming(ctx, routerService.RouteMessage)
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
