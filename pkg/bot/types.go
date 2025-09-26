package bot

// Update описывает нормализованное событие Telegram.
// TODO: сопоставить поля с реальной структурой telegram-bot-api при интеграции.
type Update struct {
	ChatID    int64  // идентификатор чата, куда нужно отвечать
	UserID    int64  // инициатор команды/сообщения
	Text      string // непустой текст сообщения
	Command   string // название команды без слеша (для /list -> "list")
	Args      string // аргументы команды без имени
	IsCommand bool   // true, если сообщение — команда Telegram
	Raw       any    // оригинальный апдейт для тонкой обработки (подключить позже)
	// TODO: заменить тип Raw на конкретный *tgbotapi.Update, когда библиотека готова
}
