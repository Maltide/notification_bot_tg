package store

import "github.com/Maltide/notification_bot_tg/pkg/types"

func (ms *MemoryStore) NextTask() (types.Task, error) {
	// TODO найти и вернуть наиближайшую задачу
	return types.Task{}, nil
}
