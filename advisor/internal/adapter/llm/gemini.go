package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lucasschilin/s5n/advisor/internal/domain"
	"google.golang.org/genai"
)

type GeminiAdapter struct {
	client *genai.Client
	model  string
}

func NewGeminiAdapter(ctx context.Context, apiKey, model string) (*GeminiAdapter, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("fail to create Gemini client: %w", err)
	}

	return &GeminiAdapter{
		client: client,
		model:  model,
	}, nil
}

func (a *GeminiAdapter) ClassifyIntent(
	ctx context.Context, userText, fallbackMessageQueueName string,
) (*domain.ClassifiedIntent, error) {
	preprompt := domain.RouterPrompt(fallbackMessageQueueName)
	prompt := fmt.Sprintf(
		"%s\n\nData/Hora Atual: %s\nMensagem do usuário: %s",
		preprompt,
		time.Now().UTC().Format(time.RFC3339),
		userText,
	)

	fmt.Println(prompt)

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"target_queue": {Type: genai.TypeString},
				"action":       {Type: genai.TypeString},
				"confidence":   {Type: genai.TypeNumber},
				"parameters":   {Type: genai.TypeObject},
			},
			Required: []string{"target_queue", "action", "confidence", "parameters"},
		},
		Temperature: genai.Ptr(float32(0.1)),
	}

	resp, err := a.client.Models.GenerateContent(ctx, a.model, genai.Text(prompt), config)
	if err != nil {
		return nil, fmt.Errorf("fail to call Gemini API: %w", err)
	}

	rawText := resp.Text()
	if rawText == "" {
		return nil, fmt.Errorf("empty response from Gemini API")
	}

	var intent domain.ClassifiedIntent
	if err := json.Unmarshal([]byte(rawText), &intent); err != nil {
		return nil, fmt.Errorf("fail to unmarshal response (%s): %w", rawText, err)
	}

	return &intent, nil
}
