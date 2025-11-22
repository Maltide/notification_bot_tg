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

	nearTask := types.Task{}
	now := time.Now()

	for _, val := range ms.Data {
		for _, task := range val {
			if nearTask.UserID == 0 {
				nearTask = task
				continue
			}
			if task.DueAt.IsZero() {
				ms.Logger.Warn("find task with zero time")
				continue
			}
			if task.DueAt.Before(now) {
				ms.Logger.Warn("find task in the past")
				return task, nil
			}
			if task.DueAt.Before(nearTask.DueAt) {
				nearTask = task
				continue
			}
		}
	}

	ms.Logger.Infof("near task to send:%v\n", nearTask)
	if nearTask == (types.Task{}) {
		return types.Task{}, ErrNoTasks
	}
	return nearTask, nil
}
