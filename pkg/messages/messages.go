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

func MsgHelp() string {
	return `Доступные команды:
/add <дата> <время> <текст> - Добавляет новую заметку(пример: /add 12.05.2025 15:05 Купить продукты)
/list - Покажет Вам текущие заметки и её ID(по нему можно удалить заметку)
/delete <id> - Удаляет заметку по ID (пример: вызвали /list - увидели какая цифра ID у заметки и потом эту цифру пишите в команде /delete 1, где 1 - ID заметки)
/help - Показать все доступные команды`
}
