package store

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

func TestMemoryStore_Race(t *testing.T) {
	s := &MemoryStore{
		Data:   make(map[int64]map[int64]types.Task),
		Logger: zap.NewNop().Sugar(),
		nextID: 1,
	}

	var wg sync.WaitGroup
	ctx := context.Background()

	// 1. Параллельные CreateTask
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(u int64) {
			defer wg.Done()
			_, _ = s.CreateTask(ctx, types.Task{
				UserID: u,
				Text:   "race",
				DueAt:  time.Now().Add(1 * time.Hour),
			})
		}(int64(i % 5))
	}

	// 2. Параллельные ListTasks
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(u int64) {
			defer wg.Done()
			_, _ = s.ListTasks(ctx, u)
		}(int64(i % 5))
	}

	// 3. Параллельные DeleteTask
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(u int64) {
			defer wg.Done()
			// пробуем удалить всё подряд — ошибки игнорируем
			_ = s.DeleteTask(ctx, u, int64(i))
		}(int64(i % 5))
	}

	wg.Wait()
}
