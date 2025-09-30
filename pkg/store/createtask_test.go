package store

import (
	"context"
	"strings"
	"testing"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

func TestCreateTask(t *testing.T) {
	store := &MemoryStore{Data: make(map[int64]map[int64]types.Task), Logger: zap.NewNop().Sugar(), nextID: 1}

	task := types.Task{UserID: 1, Text: "test"}

	created, err := store.CreateTask(context.Background(), task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created.ID == 0 {
		t.Errorf("expected non-zero ID, got %d", created.ID)
	}

	if created.Text != task.Text {
		t.Errorf("expected text %q, got %q", task.Text, created.Text)
	}
	if strings.TrimSpace(created.Text) == "" {
		t.Error("No text")
	}

	tasks, _ := store.ListTasks(context.Background(), task.UserID)
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}

	t.Run("empty text", func(t *testing.T) {
		store := &MemoryStore{Data: make(map[int64]map[int64]types.Task), Logger: zap.NewNop().Sugar(), nextID: 1}

		snap := store.nextID

		_, err := store.CreateTask(context.Background(), types.Task{UserID: 1, Text: ""})
		if err == nil {
			t.Fatal("want error for empty text")
		}

		tasks, _ := store.ListTasks(context.Background(), 1)
		if len(tasks) != 0 {
			t.Fatalf("want 0 tasks, got %d", len(tasks))
		}
		if store.nextID != snap {
			t.Fatalf("nextID changed: %d -> %d", snap, store.nextID)
		}
	})
}
