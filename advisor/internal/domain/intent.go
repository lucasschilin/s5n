package domain

import "encoding/json"

const (
	QueueReminder = "reminder"
	QueueUnknown  = "unknown_queue"
)

type ClassifiedIntent struct {
	TargetQueue string          `json:"target_queue"`
	Action      string          `json:"action"`
	Confidence  float64         `json:"confidence"`
	Parameters  json.RawMessage `json:"parameters"`
}

type RoutedMessage struct {
	OriginalMessage RawIncomingMessage `json:"original_message"`
	Intent          ClassifiedIntent   `json:"intent"`
}
