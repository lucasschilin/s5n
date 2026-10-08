package services

import (
	"context"
	"time"

	"github.com/lucasschilin/s5n/reminder/internal/domain"
	"github.com/lucasschilin/s5n/reminder/internal/ports"
)

type ReminderService struct {
	repo ports.ReminderRepository
}

func NewReminderService(repo ports.ReminderRepository) ports.ReminderService {
	return &ReminderService{repo: repo}
}

func (s *ReminderService) ProcessQueueMessage(ctx context.Context, payload domain.ReminderQueuePayload) error {
	now := time.Now().UTC()

	reminder := &domain.Reminder{
		UserID:      payload.OriginalMessage.UserID,
		Title:       payload.Parameters.Title,
		Description: payload.Parameters.Description,
		RemindAt:    payload.Parameters.RemindAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.repo.Save(ctx, reminder)
}
