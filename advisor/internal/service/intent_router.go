package service

import (
	"context"
	"fmt"
	"log"

	"github.com/lucasschilin/s5n/advisor/internal/domain"
)

type IntentRouterService struct {
	llm      domain.Classifier
	producer domain.QueueProducer
}

func NewIntentRouterService(llm domain.Classifier, producer domain.QueueProducer) *IntentRouterService {
	return &IntentRouterService{
		llm:      llm,
		producer: producer,
	}
}

func (s *IntentRouterService) RouteMessage(ctx context.Context, msg domain.RawIncomingMessage) error {
	intent, err := s.llm.ClassifyIntent(ctx, msg.Body)
	if err != nil {
		return fmt.Errorf("fail	on LLM classification: %w", err)
	}

	log.Printf("🎯 Indentified Intent: Queue=[%s] Action=[%s] Confidence=[%.2f]",
		intent.TargetQueue, intent.Action, intent.Confidence)

	routed := domain.RoutedMessage{
		OriginalMessage: msg,
		Intent:          *intent,
	}

	if err := s.producer.EnqueueToQueue(ctx, intent.TargetQueue, routed); err != nil {
		return fmt.Errorf("falha ao publicar na fila destino %s: %w", intent.TargetQueue, err)
	}

	return nil
}
