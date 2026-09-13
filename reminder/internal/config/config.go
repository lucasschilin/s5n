package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQConnURL string
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("Any .env file finded.")
	}

	AppConfig = &Config{
		RabbitMQConnURL: os.Getenv("QUEUE_RABBITMQ_CONN_URL"),
	}
}
