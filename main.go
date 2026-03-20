package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	botpkg "github.com/Maltide/notification_bot_tg/pkg/bot"
	"github.com/Maltide/notification_bot_tg/pkg/config"
	"github.com/Maltide/notification_bot_tg/pkg/db"
	helperpkg "github.com/Maltide/notification_bot_tg/pkg/helpers"
	"github.com/Maltide/notification_bot_tg/pkg/logger"
	parserpkg "github.com/Maltide/notification_bot_tg/pkg/parser"
	"github.com/Maltide/notification_bot_tg/pkg/scheduler"
	"github.com/Maltide/notification_bot_tg/pkg/store"
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

	ctx := context.Background()

	db := db.DBInit(log, cfg, ctx)

	postgStore := store.NewPostgresStore(db, log)

	// no explicit signal handling; rely on Docker/OS to stop the process

	var wg sync.WaitGroup

	api, err := tgbotapi.NewBotAPI(cfg.TGToken)
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	parser := parserpkg.NewTimeParser()

	// create handler without scheduler to avoid constructor cycle
	handler := helperpkg.NewHandler(postgStore, nil, log, parser)

	// create bot with handler
	bot := botpkg.NewBot(api, log, handler, &cfg)

	// create scheduler and pass the bot as Notifier
	sched := scheduler.NewTimerScheduler(postgStore, bot, log)

	// attach scheduler to handler now that it exists
	handler.SetScheduler(sched)

	wg.Add(1)
	go func() {
		defer wg.Done()
		sched.Start(ctx)
	}()

	// Важно: после рестарта задачи уже лежат в БД. Refresh заставит scheduler выбрать ближайшую.
	sched.Refresh()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bot.Start(ctx)
	}()

	// block until process is stopped externally (e.g., Docker sends SIGTERM)
	select {}
}
