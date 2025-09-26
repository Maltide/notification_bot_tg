package bot

import (
	"context"

	"go.uber.org/zap"
)

// Poller описывает источник обновлений (long polling, вебхук и т.п.).
// TODO: привязать к telegram-bot-api и передавать параметры offset/timeout.
type Poller interface {
	Updates(ctx context.Context) (<-chan Update, error)
}

// Bot объединяет источник обновлений и маршрутизатор.
// TODO: расширить зависимостями (formatter, error handler) при необходимости.
type Bot struct {
	Poller Poller
	Router Router
	Logger *zap.SugaredLogger
}

// Start запускает бесконечный цикл обработки апдейтов Telegram.
func (b *Bot) Start(ctx context.Context) error {
	// TODO: получить канал обновлений у Poller, запускать обработку в отдельной горутине.
	return nil
}

// dispatch передаёт обновление в маршрутизатор в зависимости от типа сообщения.
func (b *Bot) dispatch(ctx context.Context, upd Update) error {
	// TODO: вызвать HandleMessage или HandleCommand в зависимости от upd.IsCommand.
	return nil
}
