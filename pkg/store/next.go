package store

import (
	"context"
	"fmt"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
)

func (ms *MemoryStore) NextTask(ctx context.Context) (types.Task, error) {
	if ms.Data == nil {
		return types.Task{}, fmt.Errorf("no any cases at all")
	}

	copied := false
	nearTask := types.Task{}

	for _, val := range ms.Data {
		now := time.Now()
		for _, task := range val {
			if task.DueAt.Before(now) {
				ms.Logger.Warn("find task in the past")
				continue
			}
			if task.DueAt.IsZero() {
				ms.Logger.Warn("find task with zero time")
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
	ms.Logger.Infof("near task to send:%v\n", nearTask)
	return nearTask, nil
}
