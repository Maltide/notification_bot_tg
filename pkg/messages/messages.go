package messages

import (
	"fmt"
	"strings"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
)

// MsgDelete returns the user-facing text for a successful delete.
func MsgDelete() string {
	return "Заметка успешно удалена"
}

// MsgAdd returns the user-facing text for a newly created task.
func MsgAdd(createdTaskID int64, due time.Time) string {
	return fmt.Sprintf("Заметка добавлена! ID: %d | Придет: %s в %s\n", createdTaskID, due.Format("02.01.06"), due.Format("15:04"))
}

// MsgList formats a list of tasks for displaying to the user.
func MsgList(tasks []types.Task) string {
	var sb strings.Builder
	sb.WriteString("Ваши заметки:\n")
	for _, task := range tasks {
		sb.WriteString(fmt.Sprintf("ID: %d | %s | Придет: %s в %s\n",
			task.UserTaskID, task.Text, task.DueAt.Format("02.01.06"), task.DueAt.Format("15:04")))
	}
	return sb.String()
}

// MsgHelp returns help text listing supported bot commands.
func MsgHelp() string {
	return `Доступные команды:

/add DD.MM.YY hh:mm Текст - Добавляет новую заметку.
Пример: /add 12.05.25 15:05 Купить продукты
(можно без заглавной буквы)

/list - Покажет Вам текущие заметки и её ID.
(по ID можно удалить заметку)

/delete номер ID - Удаляет заметку по ID .
Пример: вызвали /list - увидели какая цифра ID у заметки и потом эту цифру пишите в команде /delete 1, где 1 - ID заметки

/help - Показать все доступные команды.`
}
