package store

import (
	"context"
	"testing"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name    string
		in      types.Task
		wantErr bool
	}{
		{"store is not empty after create", types.Task{UserID: 1, Text: "not empty", DueAt: time.Time{}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &MemoryStore{
				Data:   make(map[int64]map[int64]types.Task),
				Logger: zap.NewNop().Sugar(),
				nextID: 1,
			}

			got, err := store.CreateTask(context.Background(), tt.in)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got.DueAt != tt.in.DueAt {
				t.Fatalf("dueAt mismatch")
			}
			if got.ID == 0 {
				t.Fatalf("want non-zero ID")
			}
			if got.Text != tt.in.Text {
				t.Fatalf("text mismatch: %q vs %q", got.Text, tt.in.Text)
			}
			_, ok := store.Data[tt.in.UserID]
			if !ok {
				t.Fatalf("task was not added")
			}
			taskinStore := store.Data[got.UserID][got.ID]
			if taskinStore.Text != tt.in.Text {
				t.Fatalf("text was not added to store")
			}
			if taskinStore.DueAt != tt.in.DueAt {
				t.Fatalf("dueAt was not added to store")
			}
		})
	}
}
