package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lucasschilin/s5n/advisor/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMQProducer(amqpURL, queueName string) (*RabbitMQProducer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel in RabbitMQ: %w", err)
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

	return &RabbitMQProducer{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
	}, nil
}

func (r *RabbitMQProducer) EnqueueToQueue(ctx context.Context, targetQueue string, payload any) error {
	_, err := r.channel.QueueDeclare(
		targetQueue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare target queue %s: %w", targetQueue, err)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("fail to marshal payload: %w", err)
	}

	err = r.channel.PublishWithContext(
		ctx,
		"",          // default exchange
		targetQueue, // routing key = nome da fila destino
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("fail to publish payload to queue %s: %w", targetQueue, err)
	}

	return nil
}

func (r *RabbitMQProducer) EnqueueIncomingMessage(ctx context.Context, message *domain.RawIncomingMessage) error {
	return r.EnqueueToQueue(ctx, r.queueName, message)
}

func (r *RabbitMQProducer) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
