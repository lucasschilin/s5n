package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandlerFunc[T any] func(ctx context.Context, msg T) error

type RabbitMQConsumer[T any] struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMQConsumer[T any](amqpURL, queueName string) (*RabbitMQConsumer[T], error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	if queueName != "" {
		_, err = ch.QueueDeclare(
			queueName,
			true,  // durable
			false, // auto-delete when there are no consumers
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}
	}

	return &RabbitMQConsumer[T]{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
	}, nil
}

func (c *RabbitMQConsumer[T]) StartConsuming(ctx context.Context, handler MessageHandlerFunc[T]) error {
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
		var msg T
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			log.Printf("[ERRO] Corrupted payload, sending NACK without requeue: %v", err)
			d.Nack(false, false)
			continue
		}

		if err := handler(ctx, msg); err != nil {
			log.Printf("[ERRO PROCESSAMENTO] re-queueing message on queue [%s]: %v", c.queueName, err)
			d.Nack(false, true)
			continue
		}

		d.Ack(false)
		log.Printf("✅ Message processed and ACKed on queue %s", c.queueName)
	}

	return nil
}

func (r *RabbitMQConsumer[T]) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
