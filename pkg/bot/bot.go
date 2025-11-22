package bot

import (
	"context"
	"strings"

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
func NewBot(api *tgbotapi.BotAPI, logger *zap.SugaredLogger, cmdHandler *helperpkg.Handler,
	notifyCh <-chan types.Task) *Bot {
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
	b.sendMessage(upd.Message.Chat.ID, "Use /help")
	// TODO: здесь будет обработка обычных сообщений (создание напоминаний).
}

func (b *Bot) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	command := "/" + msg.Command()
	commandArgs := msg.CommandArguments()
	args := strings.Fields(commandArgs)
	response := b.cmdHandler.HandleCommand(ctx, msg.From.ID, msg.Chat.ID, command, args)
	reply := tgbotapi.NewMessage(msg.Chat.ID, response)
	b.sendMessage(reply.ChatID, reply.Text)
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
			b.sendMessage(task.ChatID, task.Text)
		}
	}
}

func (b *Bot) sendMessage(chatID int64, text string) error {
	message := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(message); err != nil {
		b.logger.Errorf("failed to send message: %v", err)
	}
	return nil
}
