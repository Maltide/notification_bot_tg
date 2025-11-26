package store

import (
	"context"
	"testing"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

func TestDeleteTask(t *testing.T) {
	newStore := func() *MemoryStore {
		return &MemoryStore{Data: make(map[int64]map[int64]types.Task), Logger: zap.NewNop().Sugar(), nextID: 1}
	}

	type args struct{ userID, userTaskID int64 }

	tests := []struct {
		name    string
		seed    func(s *MemoryStore) // чем наполняем стор перед вызовом
		args    args
		wantErr bool
	}{
		{
			name: "ok: delete existing",
			seed: func(s *MemoryStore) {
				u := int64(1)
				tk := types.Task{ID: 10, UserID: u, Text: "del", UserTaskID: 1}
				s.Data[u] = map[int64]types.Task{10: tk}
			},
			args:    args{userID: 1, userTaskID: 1},
			wantErr: false,
		},
		{
			name:    "err: user not found",
			seed:    func(s *MemoryStore) {}, // пусто
			args:    args{userID: 999, userTaskID: 1},
			wantErr: true,
		},
		{
			name:    "err: task not found",
			seed:    func(s *MemoryStore) { s.Data[1] = map[int64]types.Task{} },
			args:    args{userID: 1, userTaskID: 999},
			wantErr: true,
		},
		{
			name: "err: second delete",
			seed: func(s *MemoryStore) {
				u := int64(2)
				tk := types.Task{ID: 20, UserID: u, Text: "once", UserTaskID: 1}
				s.Data[u] = map[int64]types.Task{20: tk}
			},
			args:    args{userID: 2, userTaskID: 1},
			wantErr: false, // первый вызов
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStore()
			tt.seed(s)

			err := s.DeleteTask(context.Background(), tt.args.userID, tt.args.userTaskID)
			if tt.name == "err: second delete" {
				if err != nil {
					t.Fatalf("unexpected err on first delete: %v", err)
				}
				err = s.DeleteTask(context.Background(), tt.args.userID, tt.args.userTaskID) // второй раз
				if err == nil {
					t.Fatalf("expected error on second delete, got nil")
				}
				return
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
			if !tt.wantErr {
				if _, ok := s.Data[tt.args.userID][tt.args.userTaskID]; ok {
					t.Fatalf("task still in store after delete")
				}
			}
		})
	}
}
