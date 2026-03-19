package types

import (
	"context"
	"time"
)

// Task represents a scheduled reminder.
type Task struct {
	ID, UserID, ChatID, UserTaskID int64
	Text                           string
	DueAt                          time.Time
}

// Store is the persistence interface for tasks.
type Store interface {
	CreateTask(ctx context.Context, t Task) (Task, error)
	ListTasks(ctx context.Context, userID int64) ([]Task, error)
	DeleteTask(ctx context.Context, userID, id int64) error
	NextTask(ctx context.Context) (Task, error)
}

// Parser parses user input for commands (e.g. /add).
type Parser interface {
	ParseAddTask(userargs []string) (due time.Time, text string, err error)
}

// Notifier delivers a task notification to the end user (e.g. via Telegram).
type Notifier interface {
	Notify(ctx context.Context, task Task) error
}
