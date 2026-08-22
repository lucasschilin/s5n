package domain

import "context"

type Classifier interface {
	ClassifyIntent(ctx context.Context, userText, fallbackMessageQueueName string) (*ClassifiedIntent, error)
}
