package store

import (
	"context"
	"testing"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

func TestDeleteTask(t *testing.T) {
	store := &MemoryStore{
		Data:   make(map[int64]map[int64]types.Task),
		Logger: zap.NewNop().Sugar(),
		nextID: 1,
	}

	userID := int64(1)

	task := types.Task{
		ID:     1,
		UserID: userID,
		Text:   "ready to delete this note",
	}

	store.Data[userID] = map[int64]types.Task{
		task.ID: task,
	}

	err := store.DeleteTask(context.Background(), userID, task.ID)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if _, ok := store.Data[userID][task.ID]; ok {
		t.Error("task not deleted, still in store")
	}

	err = store.DeleteTask(context.Background(), userID, task.ID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
