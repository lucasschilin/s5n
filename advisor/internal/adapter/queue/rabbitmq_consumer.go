package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/lucasschilin/s5n/advisor/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandlerFunc func(ctx context.Context, msg domain.RawIncomingMessage) error

type RabbitMQConsumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMQConsumer(amqpURL, queueName string) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
	}, nil
}

func (c *RabbitMQConsumer) StartConsuming(ctx context.Context, handler MessageHandlerFunc) error {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("fail to consume queue %s: %w", c.queueName, err)
	}

	log.Printf("👷 Consumer started for queue: %s", c.queueName)

	for d := range msgs {
		var rawMsg domain.RawIncomingMessage
		if err := json.Unmarshal(d.Body, &rawMsg); err != nil {
			log.Printf("[ERRO] Corrupted payload, sending NACK without requeue: %v", err)
			d.Nack(false, false)
			continue
		}

		if err := handler(ctx, rawMsg); err != nil {
			log.Printf("[ERRO PROCESSAMENTO] re-queueing message %s: %v", rawMsg.MessageID, err)
			d.Nack(false, true)
			continue
		}

		d.Ack(false)
		log.Printf("✅ Message %s processed and removed from queue %s", rawMsg.MessageID, c.queueName)
	}

	return nil
}
