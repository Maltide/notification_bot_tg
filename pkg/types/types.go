package types

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Task struct {
	ID, UserID, ChatID int64
	Text               string
	DueAt              time.Time
	cancel             func()
}
type MemoryStore struct {
	mu     sync.RWMutex
	nextID atomic.Int64
	data   map[int64]map[int64]Task
}

type Store interface {
	CreateTask(ctx context.Context, t Task) (Task, error)
	ListTasks(ctx context.Context, userID int64, limit int) ([]Task, error)
	DeleteTask(ctx context.Context, userID, id int64) error
}
