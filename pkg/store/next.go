package store

import (
	"context"

	"github.com/Maltide/notification_bot_tg/pkg/types"
)

func (ms *MemoryStore) NextTask(ctx context.Context) (types.Task, error) {
	// TODO найти и вернуть наиближайшую задачу
	return types.Task{}, nil
}
