package store

import (
	"context"
	"testing"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

func TestListTask(t *testing.T) {
	ctx := context.Background()

	newStore := func() *MemoryStore {
		return &MemoryStore{
			Data:   make(map[int64]map[int64]types.Task),
			Logger: zap.NewNop().Sugar(),
			nextID: 1}
	}

	tests := []struct {
		name    string
		seed    func(*MemoryStore)
		userID  int64
		wantlen int
		wantErr bool
	}{
		{
			name:    "no tasks",
			seed:    func(*MemoryStore) {},
			userID:  1,
			wantlen: 0,
			wantErr: true,
		},
		{
			name: "correct len(slice) of tasks",
			seed: func(s *MemoryStore) {
				s.Data[1] = map[int64]types.Task{
					10: {ID: 10, UserID: 1, Text: "task"},
				}
			},
			userID:  1,
			wantlen: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStore()
			tt.seed(s)

			tasks, err := s.ListTasks(ctx, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v, got error: %v", tt.wantErr, err)
			}
			if len(tasks) != tt.wantlen {
				t.Fatalf("len mismatch: have %v, need %v", len(tasks), tt.wantlen)
			}

			for _, task := range tasks {
				taskInStore, ok := s.Data[tt.userID][task.ID]
				if !ok {
					t.Fatalf("task with ID %v not found in store", task.ID)
				}
				if taskInStore.Text != task.Text {
					t.Fatalf("want text=%v, got text=%v", task.Text, taskInStore.Text)
				}
				if taskInStore.DueAt != task.DueAt {
					t.Fatalf("want dueAt=%v, got dueAt=%v", task.DueAt, taskInStore.DueAt)
				}
			}
		})
	}
}
