package ports

import (
	"context"

	"github.com/lucasschilin/s5n/reminder/internal/domain"
)

type ReminderService interface {
	ProcessQueueMessage(ctx context.Context, payload domain.ReminderQueuePayload) error
}
