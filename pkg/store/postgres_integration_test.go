package store

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set; skipping Postgres integration tests")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
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

func TestPostgresStore_CreateListDelete_ShiftIDs(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	truncateTasks(t, db)

	ps := NewPostgresStore(db, zap.NewNop().Sugar())
	ctx := context.Background()

	due := time.Now().Add(1 * time.Hour)

	// create 3 tasks
	for i := 0; i < 3; i++ {
		_, err := ps.CreateTask(ctx, types.Task{UserID: 1, ChatID: 100, Text: "t", DueAt: due.Add(time.Duration(i) * time.Minute)})
		if err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	list, err := ps.ListTasks(ctx, 1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("want 3 tasks, got %d", len(list))
	}
	if list[0].UserTaskID != 1 || list[1].UserTaskID != 2 || list[2].UserTaskID != 3 {
		t.Fatalf("unexpected user_task_id sequence: %d %d %d", list[0].UserTaskID, list[1].UserTaskID, list[2].UserTaskID)
	}

	// delete middle and ensure shift
	if err := ps.DeleteTask(ctx, 1, 2); err != nil {
		t.Fatalf("delete: %v", err)
	}

	list2, err := ps.ListTasks(ctx, 1)
	if err != nil {
		t.Fatalf("list2: %v", err)
	}
	if len(list2) != 2 {
		t.Fatalf("want 2 tasks, got %d", len(list2))
	}
	if list2[0].UserTaskID != 1 || list2[1].UserTaskID != 2 {
		t.Fatalf("shift failed, got ids: %d %d", list2[0].UserTaskID, list2[1].UserTaskID)
	}
}

func TestPostgresStore_NextTask(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	truncateTasks(t, db)

	ps := NewPostgresStore(db, zap.NewNop().Sugar())
	ctx := context.Background()

	_, err := ps.CreateTask(ctx, types.Task{UserID: 1, ChatID: 1, Text: "late", DueAt: time.Now().Add(2 * time.Hour)})
	if err != nil {
		t.Fatalf("create late: %v", err)
	}
	_, err = ps.CreateTask(ctx, types.Task{UserID: 2, ChatID: 2, Text: "soon", DueAt: time.Now().Add(10 * time.Minute)})
	if err != nil {
		t.Fatalf("create soon: %v", err)
	}

	next, err := ps.NextTask(ctx)
	if err != nil {
		t.Fatalf("next: %v", err)
	}
	if next.Text != "soon" {
		t.Fatalf("want soon, got %q", next.Text)
	}
}

func TestPostgresStore_NoTasks(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	truncateTasks(t, db)

	ps := NewPostgresStore(db, zap.NewNop().Sugar())
	_, err := ps.NextTask(context.Background())
	if err == nil {
		t.Fatalf("want error")
	}
	if err != ErrNoTasks {
		t.Fatalf("want ErrNoTasks, got %v", err)
	}
}
