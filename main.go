package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	botpkg "github.com/Maltide/notification_bot_tg/pkg/bot"
	"github.com/Maltide/notification_bot_tg/pkg/config"
	helperpkg "github.com/Maltide/notification_bot_tg/pkg/helpers"
	"github.com/Maltide/notification_bot_tg/pkg/logger"
	parserpkg "github.com/Maltide/notification_bot_tg/pkg/parser"
	"github.com/Maltide/notification_bot_tg/pkg/scheduler"
	"github.com/Maltide/notification_bot_tg/pkg/store"
	"github.com/Maltide/notification_bot_tg/pkg/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	notifyCh := make(chan types.Task, 1)

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "main: config:", err)
		os.Exit(1)
	}

	log, err := logger.Logger(cfg.LogLevel)
	if err != nil {
		os.Exit(1)
	}
	log.Info("main: logger started")

	api, err := tgbotapi.NewBotAPI(cfg.TGToken)
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	ms := store.NewMemoryStore(log)

	notifyUser := func(ctx context.Context, task types.Task) error {
		log.Infof("notify task: %v", task)
		notifyCh <- task
		return nil
	}

	scheduler := scheduler.NewTimerScheduler(ms, notifyUser, log)

	wg.Add(1)
	go func() {
		defer wg.Done()
		scheduler.Start(ctx)
	}()

	parser := parserpkg.NewTimeParser()

	handler := helperpkg.NewHandler(ms, scheduler, log, parser)

	bot := botpkg.NewBot(api, log, handler, notifyCh)

	wg.Add(1)
	go func() {
		defer wg.Done()
		bot.Start(ctx)
	}()

	wg.Wait()
}
