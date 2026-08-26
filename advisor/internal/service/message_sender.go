package service

import (
	"context"

	"github.com/lucasschilin/s5n/advisor/internal/domain"
)

type MessageSenderService struct {
	messenger domain.Messenger
}

func NewMessageSenderService(messenger domain.Messenger) *MessageSenderService {
	return &MessageSenderService{
		messenger: messenger,
	}
}

func (s *MessageSenderService) SendMessage(ctx context.Context, msg domain.RawOutgoingMessage) error {
	return s.messenger.SendMessage(msg.UserID, msg.Body, msg.ReplyMessageID)
}
