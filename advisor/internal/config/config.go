package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                   string
	IncomingMessagesQueueName string
	RabbitMQConnURL           string
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("Any .env file finded.")
	}

	AppConfig = &Config{
		AppPort:                   os.Getenv("APP_PORT"),
		IncomingMessagesQueueName: os.Getenv("INCOMING_MESSAGES_QUEUE_NAME"),
		RabbitMQConnURL:           os.Getenv("RABBITMQ_CONN_URL"),
	}
}
