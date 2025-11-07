package store

import (
	"context"
	"testing"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

func TestCreateTask(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		in      types.Task
		wantErr bool
	}{
		{"past", types.Task{UserID: 1, Text: "x", DueAt: now.Add(-time.Minute)}, true},
		{"future", types.Task{UserID: 1, Text: "x", DueAt: now.Add(time.Minute)}, false},
		{"empty_text", types.Task{UserID: 1, Text: "   ", DueAt: now.Add(time.Minute)}, true},
		{"current time", types.Task{UserID: 1, Text: "y", DueAt: time.Now()}, true},
		{"empty time", types.Task{UserID: 1, Text: "y", DueAt: time.Time{}}, true},
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
			if got.ID == 0 {
				t.Fatalf("want non-zero ID")
			}
			if got.Text != tt.in.Text {
				t.Fatalf("text mismatch: %q vs %q", got.Text, tt.in.Text)
			}
		})
	}
}
