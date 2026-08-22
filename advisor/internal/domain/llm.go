package domain

import "context"

type IntentClassifier interface {
	ClassifyIntent(ctx context.Context, userText, fallbackMessageQueueName string) (*ClassifiedIntent, error)
}
