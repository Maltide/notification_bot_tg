package scheduler

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/store"
	"github.com/Maltide/notification_bot_tg/pkg/types"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// fnNotifier wraps a function to implement types.Notifier for tests.
type fnNotifier struct {
	f func(context.Context, types.Task) error
}

func (n fnNotifier) Notify(ctx context.Context, task types.Task) error { return n.f(ctx, task) }

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set; skipping Postgres-backed scheduler race test")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		t.Fatalf("ping db: %v", err)
	}
	// Ensure schema exists.
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			chat_id BIGINT,
			user_task_id BIGINT NOT NULL,
			text TEXT NOT NULL,
			due_at TIMESTAMPTZ NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_tasks_user_user_task_id ON tasks(user_id, user_task_id);
		CREATE INDEX IF NOT EXISTS idx_tasks_due_at ON tasks(due_at);
	`); err != nil {
		_ = db.Close()
		t.Fatalf("ensure schema: %v", err)
	}
	return db
}

func truncateTasks(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec("TRUNCATE tasks RESTART IDENTITY"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

func TestScheduler_Race(t *testing.T) {
	logger := zap.NewNop().Sugar()

	db := openTestDB(t)
	defer db.Close()
	truncateTasks(t, db)

	st := store.NewPostgresStore(db, zap.NewNop().Sugar())

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
