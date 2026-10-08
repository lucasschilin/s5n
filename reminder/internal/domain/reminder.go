package domain

import (
	"time"
)

// Reminder represents the internal database record
type Reminder struct {
	ID          int64     `db:"id"`
	UserID      string    `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	RemindAt    time.Time `db:"remind_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// OriginalMessage captures the origin of the request (e.g., chat/telegram/whatsapp message)
type OriginalMessage struct {
	MessageID  string    `json:"message_id" db:"message_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Body       string    `json:"body" db:"body"`
	ReceivedAt time.Time `json:"received_at" db:"received_at"`
}

// CreateReminderParameters contains the processed metadata and scheduled time
type CreateReminderParameters struct {
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	RemindAt    time.Time `json:"remind_at" db:"remind_at"`
}

// IncomingPayload matches the incoming JSON structure from the queue
type ReminderQueuePayload struct {
	OriginalMessage OriginalMessage          `json:"original_message"`
	Action          string                   `json:"action"`
	Parameters      CreateReminderParameters `json:"parameters"`
}
