package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lucasschilin/s5n/reminder/internal/domain"
	"github.com/lucasschilin/s5n/reminder/internal/ports"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) ports.ReminderRepository {
	return &sqliteRepository{db: db}
}

func (r *sqliteRepository) Save(ctx context.Context, reminder *domain.Reminder) error {
	query := `
		INSERT INTO reminders (
			user_id, title, description, remind_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		reminder.UserID,
		reminder.Title,
		reminder.Description,
		reminder.RemindAt,
		reminder.CreatedAt,
		reminder.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert reminder: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	reminder.ID = id
	return nil
}
