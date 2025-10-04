package store

import (
	"context"
	"fmt"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
)

func (ms *MemoryStore) NextTask(ctx context.Context, t types.Task) (types.Task, error) {
	if ms.Data == nil {
		return types.Task{}, fmt.Errorf("No any cases at all")
	}
	tasks, ok := ms.Data[t.UserID]
	if !ok || len(tasks) == 0 {
		return types.Task{}, fmt.Errorf("No any task")
	}

	copied := false
	nearTask := types.Task{}

	for _, val := range ms.Data {
		now := time.Now()
		for _, task := range val {
			if task.DueAt.Before(now) {
				continue
			}
			if task.DueAt.IsZero() {
				continue
			}
			if !copied {
				nearTask = task
				copied = true
				continue
			}
			if task.DueAt.Before(nearTask.DueAt) || (task.DueAt.Equal(nearTask.DueAt) && task.UserID < nearTask.UserID) {
				nearTask = task
			}
		}
	}
	if !copied {
		return types.Task{}, fmt.Errorf("no valid tasks")
	}

	return nearTask, nil // как будем сравнивать этот neartask с новыми входящими заметками ?
} // O(n^2)
