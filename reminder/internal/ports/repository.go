package ports

import (
	"context"

	"github.com/lucasschilin/s5n/reminder/internal/domain"
)

type ReminderRepository interface {
	Save(ctx context.Context, reminder *domain.Reminder) error
}
