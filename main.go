package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	botpkg "github.com/Maltide/notification_bot_tg/pkg/bot"
	"github.com/Maltide/notification_bot_tg/pkg/config"
	"github.com/Maltide/notification_bot_tg/pkg/db"
	helperpkg "github.com/Maltide/notification_bot_tg/pkg/helpers"
	"github.com/Maltide/notification_bot_tg/pkg/logger"
	parserpkg "github.com/Maltide/notification_bot_tg/pkg/parser"
	"github.com/Maltide/notification_bot_tg/pkg/scheduler"
	"github.com/Maltide/notification_bot_tg/pkg/store"
	"github.com/Maltide/notification_bot_tg/pkg/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// main is the program entry point: it initializes services and starts the bot and scheduler.
func main() {
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := db.DBInit(log, cfg, ctx)

	postgStore := store.NewPostgresStore(db, log)

	endCh := make(chan os.Signal, 1)
	signal.Notify(endCh, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	notifyCh := make(chan types.Task, 1)

	api, err := tgbotapi.NewBotAPI(cfg.TGToken)
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	notifyUser := func(ctx context.Context, task types.Task) error {
		log.Infof("notify task: %v", task)
		notifyCh <- task
		return nil
	}

	scheduler := scheduler.NewTimerScheduler(postgStore, notifyUser, log)

	wg.Add(1)
	go func() {
		defer wg.Done()
		scheduler.Start(ctx)
	}()

	// Важно: после рестарта задачи уже лежат в БД. Refresh заставит scheduler выбрать ближайшую.
	scheduler.Refresh()

	parser := parserpkg.NewTimeParser()

	handler := helperpkg.NewHandler(postgStore, scheduler, log, parser)

	bot := botpkg.NewBot(api, log, handler, notifyCh, &cfg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		bot.Start(ctx)
	}()

	<-endCh
	cancel()
	scheduler.Stop()
	wg.Wait()
}
