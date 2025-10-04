package store

import (
	"testing"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
)

func TestNextTask(t *testing.T) {
	store := &MemoryStore{Data: map[int64]map[int64]types.Task{}}
	tests := []struct {
		name      string
		userinput types.Task
		wantErr   bool
	}{
		{"empty", types.Task{UserID: 1, Text: "", DueAt: time.Time{}}, true},
		{"future", types.Task{UserID: 1, Text: "x", DueAt: now.Add(time.Minute)}, false},
		{"empty_text", types.Task{UserID: 1, Text: "   ", DueAt: now.Add(time.Minute)}, true},
	}
}
