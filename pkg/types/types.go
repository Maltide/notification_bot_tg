package types

import (
	"context"
)

type Task struct {
	ID, UserID, ChatID int64
	Text               string
	// DueAt              time.Time
	// cancel             func() //????????????????????????????
}

type Store interface {
	CreateTask(ctx context.Context, t Task) (Task, error)
	ListTasks(ctx context.Context, userID int64) ([]Task, error)
	DeleteTask(ctx context.Context, userID, id int64) error
}
