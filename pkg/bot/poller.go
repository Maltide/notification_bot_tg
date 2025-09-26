package bot

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// LongPoller использует telegram-bot-api для чтения апдейтов через long polling.
// TODO: добавить поддержку оффсетов и остановку при закрытии контекста.
type LongPoller struct {
	API     *tgbotapi.BotAPI
	Timeout time.Duration
}

// Updates запускает получение апдейтов у Telegram и нормализует их в Update.
func (p *LongPoller) Updates(ctx context.Context) (<-chan Update, error) {
	// TODO: реализовать цикл общения с API и преобразование *tgbotapi.Update -> Update.
	return nil, nil
}
