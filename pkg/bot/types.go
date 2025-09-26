package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Update описывает нормализованное событие Telegram.
// TODO: уточнить набор полей после первых итераций с реальными апдейтами.
type Update struct {
	ChatID    int64            // идентификатор чата, куда нужно отвечать
	UserID    int64            // инициатор команды/сообщения
	Text      string           // непустой текст сообщения
	Command   string           // название команды без слеша (для /list -> "list")
	Args      string           // аргументы команды без имени
	IsCommand bool             // true, если сообщение — команда Telegram
	Raw       *tgbotapi.Update // оригинальный апдейт для тонкой обработки
	// TODO: добавить вспомогательные методы (Reply, Sender) при необходимости
}
