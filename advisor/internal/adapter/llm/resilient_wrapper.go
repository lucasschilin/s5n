package llm

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/lucasschilin/s5n/advisor/internal/domain"
)

type ResilientLLMWrapper struct {
	underlying   domain.IntentClassifier
	circuitMutex sync.RWMutex
	pausedUntil  time.Time
}

// NewResilientLLMWrapper encapsula qualquer adapter de LLM com resiliência
func NewResilientLLMWrapper(underlying domain.IntentClassifier) *ResilientLLMWrapper {
	return &ResilientLLMWrapper{
		underlying: underlying,
	}
}

func (w *ResilientLLMWrapper) ClassifyIntent(
	ctx context.Context, userText, fallbackMessageQueueName string,
) (*domain.ClassifiedIntent, error) {

	w.waitUntilUnpaused(ctx)

	delays := []time.Duration{
		5 * time.Second,
		10 * time.Second,
		20 * time.Second,
	}
	maxAttempts := len(delays) + 1

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		w.waitUntilUnpaused(ctx)

		intent, err := w.underlying.ClassifyIntent(ctx, userText, fallbackMessageQueueName)
		if err == nil {
			return intent, nil
		}

		lastErr = err
		log.Printf("[LLM Wrapper Warning] Attempt %d/%d as failed: %v", attempt, maxAttempts, err)

		if attempt == maxAttempts {
			break
		}

		pauseDuration := delays[attempt-1]
		w.pauseCircuit(pauseDuration)

		log.Printf("[Circuit Breaker] Pausing circuit for %v...", pauseDuration)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pauseDuration):
		}
	}

	log.Printf("[LLM Fallback] LLM retries failed (%v). Returning fallback message.", lastErr)

	return &domain.ClassifiedIntent{
		TargetQueue: fallbackMessageQueueName,
		Action:      "SEND_MESSAGE",
		Confidence:  1.0,
		Parameters:  json.RawMessage([]byte(`{"message":"Me desculpe, não vou conseguir te responder agora, tente enviar sua mensagem novamente em alguns minutos! :("}`)),
	}, nil
}

func (w *ResilientLLMWrapper) pauseCircuit(duration time.Duration) {
	w.circuitMutex.Lock()
	defer w.circuitMutex.Unlock()

	until := time.Now().Add(duration)
	if until.After(w.pausedUntil) {
		w.pausedUntil = until
	}
}

func (w *ResilientLLMWrapper) waitUntilUnpaused(ctx context.Context) {
	w.circuitMutex.RLock()
	pauseRemaining := time.Until(w.pausedUntil)
	w.circuitMutex.RUnlock()

	if pauseRemaining > 0 {
		select {
		case <-ctx.Done():
		case <-time.After(pauseRemaining):
		}
	}
}
