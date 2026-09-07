package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                   string
	IncomingMessagesQueueName string
	OutgoingMessagesQueueName string
	RabbitMQConnURL           string
	GeminiAPIKey              string
	GeminiModel               string
	TestReceiverStubMode      bool
	TestRouterStubMode        bool
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("Any .env file finded.")
	}

	testReceiverStubMode := false
	if os.Getenv("TEST_RECEIVER_STUB_MODE") == "true" {
		testReceiverStubMode = true
	}

	testRouterStubMode := false
	if os.Getenv("TEST_ROUTER_STUB_MODE") == "true" {
		testRouterStubMode = true
	}

	AppConfig = &Config{
		AppPort:                   os.Getenv("APP_PORT"),
		IncomingMessagesQueueName: os.Getenv("QUEUE_NAME_MESSAGES_INCOMING"),
		OutgoingMessagesQueueName: os.Getenv("QUEUE_NAME_MESSAGES_OUTGOING"),
		RabbitMQConnURL:           os.Getenv("QUEUE_RABBITMQ_CONN_URL"),
		GeminiAPIKey:              os.Getenv("LLM_GEMINI_API_KEY"),
		GeminiModel:               os.Getenv("LLM_GEMINI_MODEL"),
		TestReceiverStubMode:      testReceiverStubMode,
		TestRouterStubMode:        testRouterStubMode,
	}
}
