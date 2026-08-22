package domain

import "context"

// QueueProducer defines the contract for pushing messages into the queue.
type QueueProducer interface {
	EnqueueToQueue(ctx context.Context, targetQueue string, payload any) error
}
