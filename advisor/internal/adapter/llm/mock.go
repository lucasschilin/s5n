package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lucasschilin/s5n/advisor/internal/domain"
)

type LLMMockAdapter struct{}

func NewLLMMockAdapter() *LLMMockAdapter {
	return &LLMMockAdapter{}
}

func (a *LLMMockAdapter) ClassifyIntent(ctx context.Context, userText string) (*domain.ClassifiedIntent, error) {

	preprompt := domain.RouterPrompt +
		"\n\nData/Hora Atual: " + time.Now().UTC().String() +
		"\nMensagem do usuário: " + userText

	fmt.Println(preprompt)

	rawResponse := `{
		"target_queue": "reminder_agent",
		"action": "CREATE_REMINDER",
		"confidence": 0.95,
		"parameters": {
			"title": "Comprar pão",
			"scheduled_at": "18:00"
		}
	}`

	var intent domain.ClassifiedIntent
	if err := json.Unmarshal([]byte(rawResponse), &intent); err != nil {
		return nil, fmt.Errorf("fail to unmarshal response: %w", err)
	}

	if strings.TrimSpace(intent.TargetQueue) == "" {
		intent.TargetQueue = domain.QueueAgentUnknown
	}

	return &intent, nil
}
