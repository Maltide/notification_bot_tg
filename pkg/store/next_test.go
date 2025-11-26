package store

import (
	"context"
	"testing"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

func TestNextTask(t *testing.T) {
	now := time.Now()
	logger := zap.NewNop().Sugar()

	tests := []struct {
		name    string
		store   *MemoryStore
		wantErr bool
		wantUID int64
	}{
		{
			name:    "empty store",
			store:   &MemoryStore{Data: map[int64]map[int64]types.Task{}, Logger: logger},
			wantErr: true,
		},
		{
			name: "zero DueAt",
			store: &MemoryStore{Data: map[int64]map[int64]types.Task{
				1: {1: {UserID: 1, Text: "x", DueAt: time.Time{}}},
			}, Logger: logger},
			wantErr: true,
		},
		{
			name: "past only",
			store: &MemoryStore{Data: map[int64]map[int64]types.Task{
				1: {1: {UserID: 1, Text: "x", DueAt: now.Add(-time.Minute)}},
			}, Logger: logger},
			wantErr: false,
			wantUID: 1,
		},
		{
			name: "two users future (pick UID=1)",
			store: &MemoryStore{Data: map[int64]map[int64]types.Task{
				1: {1: {UserID: 1, Text: "a", DueAt: now.Add(1 * time.Minute)}},
				2: {2: {UserID: 2, Text: "b", DueAt: now.Add(2 * time.Minute)}},
			}, Logger: logger},
			wantErr: false,
			wantUID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.store.NextTask(context.Background())
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected err=%v, got %v", tt.wantErr, err)
			}

			if !tt.wantErr && got.UserID != tt.wantUID {
				t.Fatalf("expected user %d, got %d", tt.wantUID, got.UserID)
			}
		})
	}
}
