package scheduler

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/store"
	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

// fnNotifier wraps a function to implement types.Notifier for tests.
type fnNotifier struct {
	f func(context.Context, types.Task) error
}

func (n fnNotifier) Notify(ctx context.Context, task types.Task) error { return n.f(ctx, task) }

func TestScheduler_Race(t *testing.T) {
	logger := zap.NewNop().Sugar()

	st := &store.MemoryStore{
		Data:   make(map[int64]map[int64]types.Task),
		Logger: logger,
	}

	notifyFn := func(ctx context.Context, task types.Task) error {
		// пусто — нам важно только гонки поймать
		return nil
	}

	// wrap function into a types.Notifier-compatible object
	sched := NewTimerScheduler(st, fnNotifier{notifyFn}, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go sched.Start(ctx)

	var wg sync.WaitGroup

	// Создаём много параллельных задач
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			_, _ = st.CreateTask(ctx, types.Task{
				UserID: int64(i % 5),
				Text:   "race_test",
				DueAt:  time.Now().Add(time.Millisecond * time.Duration(50+i)),
			})

			// принудительно просим scheduler пересчитать
			sched.Refresh()
		}(i)
	}

	// Параллельные refresh
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sched.Refresh()
		}()
	}

	// Параллельные удаления задач
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// удаляем всё подряд — ошибки не важны
			_ = st.DeleteTask(ctx, int64(i%5), int64(i))
			sched.Refresh()
		}(i)
	}

	wg.Wait()

	cancel()     // остановить scheduler
	sched.Stop() // корректный выход из цикла
}
