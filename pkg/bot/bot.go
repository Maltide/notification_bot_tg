package bot

import (
	"context"

	helperpkg "github.com/Maltide/notification_bot_tg/pkg/helpers"
	"github.com/Maltide/notification_bot_tg/pkg/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// Bot инкапсулирует telegram-bot-api и маршрутизацию команд.
// TODO: расширить зависимостями (parser, store, scheduler) по мере реализации.
type Bot struct {
	api        *tgbotapi.BotAPI
	logger     *zap.SugaredLogger
	cmdHandler *helperpkg.Handler
	notifyCh   <-chan types.Task
}

// NewBot создаёт каркас бота с готовым клиентом.
func NewBot(api *tgbotapi.BotAPI, logger *zap.SugaredLogger, cmdHandler *helperpkg.Handler, notifyCh <-chan types.Task) *Bot {
	return &Bot{api: api, logger: logger, cmdHandler: cmdHandler, notifyCh: notifyCh}
}

// Start запускает перехват апдейтов и реагирует хотя бы на /help.
func (b *Bot) Start(ctx context.Context) error {
	// TODO: вынести конфигурацию long polling в настройки.
	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = 30

	updates := b.api.GetUpdatesChan(updateCfg)
	go b.listenNotifications(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case upd, ok := <-updates:
			if !ok {
				return nil
			}
			b.handleUpdate(ctx, upd)
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, upd tgbotapi.Update) {
	if upd.Message == nil {
		return
	}

	if upd.Message.IsCommand() {
		b.handleCommand(ctx, upd.Message)
		return
	}

	// TODO: здесь будет обработка обычных сообщений (создание напоминаний).
}

func (b *Bot) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	switch msg.Command() {
	case "help":
		b.replyHelp(ctx, msg)
	case "list":
		// TODO: добавить обработку списка когда появится store.
	case "delete":
		// TODO: добавить обработку удаления по ID.
	default:
		b.replyUnknown(ctx, msg)
	}
}

func (b *Bot) listenNotifications(ctx context.Context) {
	if b.notifyCh == nil {
		b.logger.Warn("notifyCh is nil, notifications listener is not started")
		return
	}

	for {
		select {
		case <-ctx.Done():
			b.logger.Infof("stop listenNotifications: context cancelled")
			return
		case task, ok := <-b.notifyCh:
			if !ok {
				b.logger.Infof("stop listenNotifications: notifyCh closed")
				return
			}

			message := tgbotapi.NewMessage(task.ChatID, task.Text)
			if _, err := b.api.Send(message); err != nil {
				b.logger.Errorf("failed to send notification message: %v", err)
			}
		}
	}
}

// func (b *Bot) replyHelp(ctx context.Context, msg *tgbotapi.Message) {
// 	helpText := "Привет! Отправь сообщение вида `10m позвонить маме`, чтобы создать напоминание. Доступные команды: /help, /list, /delete <id>."
// 	// TODO: вынести шаблон текста в отдельный пакет/константу и добавить локализацию.
// 	req := tgbotapi.NewMessage(msg.Chat.ID, helpText)
// 	req.ParseMode = "Markdown"
// 	if _, err := b.api.Send(req); err != nil {
// 		b.logger.Errorf("failed to send help message: %v", err)
// 	}
// }

// func (b *Bot) replyUnknown(ctx context.Context, msg *tgbotapi.Message) {
// 	text := "Я не знаю эту команду. Используй /help, чтобы увидеть доступные команды."
// 	req := tgbotapi.NewMessage(msg.Chat.ID, text)
// 	if _, err := b.api.Send(req); err != nil {
// 		b.logger.Errorf("failed to send unknown command reply: %v", err)
// 	}
// }
