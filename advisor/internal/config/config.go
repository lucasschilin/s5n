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
	TestRouterStubMode        bool
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("Any .env file finded.")
	}

	testStubMode := false
	if os.Getenv("TEST_ROUTER_STUB_MODE") == "true" {
		testStubMode = true
	}

	AppConfig = &Config{
		AppPort:                   os.Getenv("APP_PORT"),
		IncomingMessagesQueueName: os.Getenv("QUEUE_NAME_MESSAGES_INCOMING"),
		OutgoingMessagesQueueName: os.Getenv("QUEUE_NAME_MESSAGES_OUTGOING"),
		RabbitMQConnURL:           os.Getenv("QUEUE_RABBITMQ_CONN_URL"),
		GeminiAPIKey:              os.Getenv("LLM_GEMINI_API_KEY"),
		GeminiModel:               os.Getenv("LLM_GEMINI_MODEL"),
		TestRouterStubMode:        testStubMode,
	}
}
