package store

import (
	"context"

	"github.com/Maltide/notification_bot_tg/pkg/types"
)

func (ms *MemoryStore) CreateTask(ctx context.Context, t types.Task) (types.Task, error) {
	ms.mu.Lock()
	ms.data[t.UserID][ms.nextID] = t
	ms.nextID++
	ms.mu.Unlock()

	return t, nil
}
func (ms *MemoryStore) ListTasks(ctx context.Context, userID int64, limit int) ([]types.Task, error) {
	return []types.Task{}, nil
}
func (ms *MemoryStore) DeleteTask(ctx context.Context, userID, id int64) error {
	return nil
}
