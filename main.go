package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/config"
	"github.com/Maltide/notification_bot_tg/pkg/logger"
	"github.com/Maltide/notification_bot_tg/pkg/scheduler"
	"github.com/Maltide/notification_bot_tg/pkg/store"
	"github.com/Maltide/notification_bot_tg/pkg/types"
)

func main() {
	var wg sync.WaitGroup
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

	notifyUser := func(ctx context.Context, task types.Task) error {
		log.Infof("notify task: %v", task)
		return nil
	}
	ms := store.NewMemoryStore(log)
	scheduler := scheduler.NewTimerScheduler(ms, notifyUser, log)
	wg.Add(1)
	go func() {
		scheduler.Start(context.Background())
		wg.Done()
	}()

	tasks := []types.Task{ //через ParseAddTask не стал пока заморачиваться, хотел именно Start() логику првоерить
		{UserID: 1, Text: "this exemple for delete", DueAt: time.Now().Add(20 * time.Second)},
		{UserID: 1, Text: "same time exemple1", DueAt: time.Now().Add(30 * time.Second)},
		{UserID: 2, Text: "same time exemple2", DueAt: time.Now().Add(30 * time.Second)},
		{UserID: 3, Text: "this task is faster then others, but added last", DueAt: time.Now().Add(25 * time.Second)},
	}
	for _, t := range tasks {
		ms.CreateTask(context.Background(), t)
	}
	scheduler.Refresh() // отправка сигнала после каждого изменения от юзера

	tasks1, err := ms.ListTasks(context.Background(), 1)
	log.Info("User's tasks:", tasks1)
	// у меня id задачи походу начинаются с нуля, а не с единицы, что для пользователя будет выглядеть всрато
	ms.DeleteTask(context.Background(), 1, 1)
	time.Sleep(50 * time.Millisecond) // помогло избежать дата-рейс
	scheduler.Refresh()               // отправка сигнала после каждого изменения от юзера

	tasks2, err := ms.ListTasks(context.Background(), 1)
	log.Info("User's tasks:", tasks2)

	//scheduler.Stop(context.Background()) //later
	//wg.Wait()
	select {} // зациклил чтобы все тайминги прокнули, с wg.Wait() чета не помогло.
}
