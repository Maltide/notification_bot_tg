package store

import (
	"fmt"

	"github.com/Maltide/notification_bot_tg/pkg/types"
)

func (ms *MemoryStore) NextTask() (types.Task, error) {
	if len(ms.Data) == 0 {
		return types.Task{}, fmt.Errorf("Нет заметок")
	}

	copied := false
	nearTask := types.Task{}

	for _, val := range ms.Data {
		for _, task := range val {
			if !copied {
				nearTask = task
				copied = true
				continue
			}
			if task.DueAt.Before(nearTask.DueAt) {
				nearTask = task
			}
		}
	}

	return types.Task{}, nil
}
