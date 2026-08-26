package service

import (
	"context"
	"fmt"

	"github.com/lucasschilin/s5n/advisor/internal/domain"
)

type MessageSenderService struct{}

func NewMessageSenderService() *MessageSenderService {
	return &MessageSenderService{}
}

func (s *MessageSenderService) SendMessage(ctx context.Context, msg domain.RawOutgoingMessage) error {
	if msg.ReplyMessageID == "" {
		fmt.Printf("Sending message to user [%s]: '%s'\n\n", msg.UserID, msg.Body)
		return nil
	}

	fmt.Printf("Answering message [%s] to user [%s]: '%s'\n\n", msg.ReplyMessageID, msg.UserID, msg.Body)
	return nil
}
