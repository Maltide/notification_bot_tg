package store

import (
	"context"

	"github.com/Maltide/notification_bot_tg/pkg/types"
)

func (ms *MemoryStore) CreateTask(ctx context.Context, t types.Task) (types.Task, error) {
	ms.Logger.Debugf("Executing CreateTask on task %v", t)

	ms.mu.Lock()
	ms.Logger.Debugf("Mutex was locked")

	ms.Logger.Debugf("Defining taskID: %v", ms.nextID)
	t.ID = ms.nextID

	ms.Logger.Debugf("Creating map for user with tasks.")
	if ms.Data[t.UserID] == nil {
		newmap := make(map[int64]types.Task)
		ms.Data[t.UserID] = newmap
	}
	ms.Data[t.UserID][ms.nextID] = t
	ms.nextID++
	ms.Logger.Debugf("Incrementating taskID. Now is %v", ms.nextID)

	ms.mu.Unlock()
	ms.Logger.Debugf("Mutex was unlocked")

	ms.Logger.Debugf("Returning structure t for user %v, text is %s", t.UserID, t.Text)
	return t, nil
}
func (ms *MemoryStore) ListTasks(ctx context.Context, userID int64) ([]types.Task, error) {
	ms.Logger.Debugf("Give all № %v user's tasks.", userID)
	out := make([]types.Task, 0, len(ms.Data[userID]))

	ms.Logger.Debugf("Take tasks from data and fill the var out to give this var to the user")
	for _, t := range ms.Data[userID] {
		out = append(out, t)
		ms.Logger.Debugf("Task № %v just added to list", t.ID)
	}
	ms.Logger.Debugf("Done. At now, we returning final out list %v", out)
	return out, nil
}
func (ms *MemoryStore) DeleteTask(ctx context.Context, userID, id int64) error {
	ms.Logger.Debugf("Deleting №%v user's task №%v", userID, id)

	ms.mu.Lock()
	ms.Logger.Debugf("Mutex was locked")

	delete(ms.Data[userID], id)
	ms.Logger.Debugf("Task was deleted")

	ms.mu.Unlock()
	ms.Logger.Debugf("Mutex was unlocked")

	return nil
}
