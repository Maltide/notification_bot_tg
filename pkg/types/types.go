package types

import (
	"context"
	"time"
)

type Task struct {
	UserID, ChatID int64
	Text           string
	DueAt          time.Time
	cancel         func()
}

type Store interface {
	CreateTask(ctx context.Context, t Task) (Task, error)
	ListTasks(ctx context.Context, userID int64, limit int) ([]Task, error)
	DeleteTask(ctx context.Context, userID, id int64) error
}
