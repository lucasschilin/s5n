package domain

import "context"

type Classifier interface {
	ClassifyIntent(ctx context.Context, userText string) (*ClassifiedIntent, error)
}
