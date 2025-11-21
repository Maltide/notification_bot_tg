package types

import (
	"context"
	"time"
)

type Task struct {
	ID, UserID, ChatID, UserTaskID int64
	Text                           string
	DueAt                          time.Time
}

type Store interface {
	CreateTask(ctx context.Context, t Task) (Task, error)
	ListTasks(ctx context.Context, userID int64) ([]Task, error)
	DeleteTask(ctx context.Context, userID, id int64) error
	NextTask(ctx context.Context) (Task, error)
}

type Parser interface {
	ParseAddTask(userargs []string) (due time.Time, text string, err error)
}
