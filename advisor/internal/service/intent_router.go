package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/lucasschilin/s5n/advisor/internal/domain"
)

type IntentRouterService struct {
	llm                      domain.IntentClassifier
	producer                 domain.QueueProducer
	fallbackMessageQueueName string
}

func NewIntentRouterService(llm domain.IntentClassifier, producer domain.QueueProducer, fallbackMessageQueueName string) *IntentRouterService {
	return &IntentRouterService{
		llm:                      llm,
		producer:                 producer,
		fallbackMessageQueueName: fallbackMessageQueueName,
	}
}

type OutgoingParameters struct {
	Message string `json:"message"`
}

func (s *IntentRouterService) RouteMessage(ctx context.Context, msg domain.RawIncomingMessage) error {
	intent, err := s.llm.ClassifyIntent(ctx, msg.Body, s.fallbackMessageQueueName)
	if err != nil {
		return fmt.Errorf("fail	on LLM classification: %w", err)
	}

	log.Printf("🎯 Indentified Intent: Queue=[%s] Action=[%s] Confidence=[%.2f]",
		intent.TargetQueue, intent.Action, intent.Confidence)

	var payload any
	if intent.TargetQueue != s.fallbackMessageQueueName {
		payload = domain.RoutedMessage{
			OriginalMessage: msg,
			Intent:          *intent,
		}
	} else {
		var params OutgoingParameters
		if err := json.Unmarshal(intent.Parameters, &params); err != nil {
			return fmt.Errorf("failed to unmarshal outgoing parameters: %w", err)
		}

		payload = domain.RawOutgoingMessage{
			ReplyMessageID: msg.MessageID,
			UserID:         msg.UserID,
			Body:           params.Message,
		}
	}

	if err := s.producer.EnqueueToQueue(ctx, intent.TargetQueue, payload); err != nil {
		return fmt.Errorf("failed to enqueue message to queue [%s]: %w", intent.TargetQueue, err)
	}

	return nil
}
