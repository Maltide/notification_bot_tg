package store

import (
	"context"
	"fmt"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

func NewMemoryStore(logger *zap.SugaredLogger) *MemoryStore {
	return &MemoryStore{
		Logger: logger,
		Data:   make(map[int64]map[int64]types.Task),
		nextID: 1,
	}
}
func (ms *MemoryStore) CreateTask(ctx context.Context, t types.Task) (types.Task, error) {
	ms.Logger.Debugf("Executing CreateTask on task %v", t)

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.Logger.Debugf("Mutex was locked")

	ms.Logger.Debugf("Defining taskID: %v", ms.nextID)

	t.ID = ms.nextID

	ms.Logger.Debugf("Create map for user with tasks.")

	if ms.Data[t.UserID] == nil {
		ms.Data[t.UserID] = make(map[int64]types.Task)
	}

	var maxlocalID int64 = 0

	for _, existing := range ms.Data[t.UserID] {
		if existing.UserTaskID > maxlocalID {
			maxlocalID = existing.UserTaskID
		}
	}

	t.UserTaskID = maxlocalID + 1

	ms.Data[t.UserID][ms.nextID] = t

	ms.nextID++

	ms.Logger.Debugf("Incrementating taskID. Now is %v", ms.nextID)

	ms.Logger.Debugf("Mutex was unlocked")

	ms.Logger.Debugf("Returning structure t for user %v, text is %s", t.UserID, t.Text)
	return t, nil
}

func (ms *MemoryStore) ListTasks(ctx context.Context, userID int64) ([]types.Task, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if len(ms.Data[userID]) == 0 {
		return []types.Task{}, fmt.Errorf("you have no tasks")
	}

	ms.Logger.Debugf("Give all № %v user's tasks.", userID)

	out := make([]types.Task, 0, len(ms.Data[userID]))

	ms.Logger.Debugf("Take tasks from data and fill the var out to give this var to the user")

	for _, t := range ms.Data[userID] {
		out = append(out, t)
		ms.Logger.Debugf("Task № %v just added to list", t.ID)
	}

	ms.Logger.Debugf("Done. At now, we returning final out list")

	return out, nil
}

func (ms *MemoryStore) DeleteTask(ctx context.Context, userID, userTaskID int64) error {

	ms.Logger.Debugf("Deleting №%v user's task №%v", userID, userTaskID)

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.Logger.Debugf("Mutex was locked")

	tasks, ok := ms.Data[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}

	var globalKey int64 = 0
	found := false

	for gKey, task := range tasks {

		if task.UserTaskID == userTaskID {

			globalKey = gKey
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task %d not found", userTaskID)
	}

	delete(tasks, globalKey)

	if len(tasks) == 0 {
		delete(ms.Data, userID)
		return nil
	}

	ms.Logger.Debugf("Task was deleted")

	for gKey, task := range tasks {
		if task.UserTaskID > userTaskID {
			task.UserTaskID--
			tasks[gKey] = task
		}
	}
	ms.Logger.Debugf("Mutex was unlocked")

	return nil
}
