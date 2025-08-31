package store

import (
	"context"

	"github.com/Maltide/notification_bot_tg/pkg/types"
)

func CreateTask(ctx context.Context, t types.Task) (types.Task, error)
func ListTasks(ctx context.Context, userID int64, limit int) ([]types.Task, error)
func DeleteTask(ctx context.Context, userID, id int64) error
