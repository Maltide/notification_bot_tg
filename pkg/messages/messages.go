package messages

import (
	"fmt"
	"strings"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
)

func MsgDelete() string {
	return "Task was deleted"
}

func MsgAdd(createdTaskID int64, due time.Time) string {
	return fmt.Sprintf("Task added successfully! ID: %d, Will notify you at %s", createdTaskID, due)
}

func MsgList(tasks []types.Task) string {
	var sb strings.Builder
	sb.WriteString("Your tasks:\n")
	for _, task := range tasks {
		timeUntil := time.Until(task.DueAt)
		sb.WriteString(fmt.Sprintf("ID: %d | %s | Due in: %s\n", task.ID, task.Text, timeUntil.Round(time.Second)))
	}
	return sb.String()
}
