package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
)

func (ms *MemoryStore) CreateTask(ctx context.Context, t types.Task) (types.Task, error) {
	if strings.TrimSpace(t.Text) == "" {
		return types.Task{}, fmt.Errorf("No text")
	}

	if t.DueAt.Before(time.Now().Local()) || t.DueAt.IsZero() {
		return types.Task{}, fmt.Errorf("Past time")
	}

	ms.Logger.Debugf("Executing CreateTask on task %v", t)

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.Logger.Debugf("Mutex was locked")

	ms.Logger.Debugf("Defining taskID: %v", ms.nextID)
	t.ID = ms.nextID

	ms.Logger.Debugf("Create map for user with tasks.")
	if ms.Data[t.UserID] == nil {
		newmap := make(map[int64]types.Task)
		ms.Data[t.UserID] = newmap
	}
	ms.Data[t.UserID][ms.nextID] = t
	ms.nextID++
	ms.Logger.Debugf("Incrementating taskID. Now is %v", ms.nextID)

	// ms.mu.Unlock() // нужно ли мьютекс перенести в defer чтобы он return захватывал?
	ms.Logger.Debugf("Mutex was unlocked")

	ms.Logger.Debugf("Returning structure t for user %v, text is %s", t.UserID, t.Text)
	return t, nil
}
func (ms *MemoryStore) ListTasks(ctx context.Context, userID int64) ([]types.Task, error) {
	ms.Logger.Debugf("Give all № %v user's tasks.", userID)
	out := make([]types.Task, 0, len(ms.Data[userID]))

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.Logger.Debugf("Take tasks from data and fill the var out to give this var to the user")
	for _, t := range ms.Data[userID] {
		out = append(out, t)
		ms.Logger.Debugf("Task № %v just added to list", t.ID)
	}
	ms.Logger.Debugf("Done. At now, we returning final out list")
	return out, nil
}
func (ms *MemoryStore) DeleteTask(ctx context.Context, userID, id int64) error {
	ms.Logger.Debugf("Deleting №%v user's task №%v", userID, id)

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.Logger.Debugf("Mutex was locked")

	tasks, ok := ms.Data[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}
	if _, ok := tasks[id]; !ok {
		return fmt.Errorf("task %d not found", id)
	}

	delete(tasks, id)
	if len(tasks) == 0 {
		delete(ms.Data, userID)
	}
	ms.Logger.Debugf("Task was deleted")

	ms.Logger.Debugf("Mutex was unlocked")

	return nil
}
