package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Maltide/notification_bot_tg/pkg/config"
	"github.com/Maltide/notification_bot_tg/pkg/logger"
	"github.com/Maltide/notification_bot_tg/pkg/store"
	"github.com/Maltide/notification_bot_tg/pkg/types"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "config:", err)
		os.Exit(1)
	}
	log, err := logger.Logger(cfg.LogLevel)
	if err != nil {
		os.Exit(1)
	}
	log.Info("start")

	ms := store.MemoryStore{Logger: log, Data: make(map[int64]map[int64]types.Task)}
	tasks := []types.Task{
		{UserID: 1, Text: "divan"},
		{UserID: 1, Text: "divanqjikefrhoiqefr"},
		{UserID: 2, Text: "divan2"},
		{UserID: 3, Text: "divan3"},
	}
	for _, t := range tasks {
		ms.CreateTask(context.Background(), t)
	}

	tasks1, err := ms.ListTasks(context.Background(), 1)
	log.Info("User's tasks:", tasks1)

	ms.DeleteTask(context.Background(), 1, 1)
	tasks2, err := ms.ListTasks(context.Background(), 1)
	log.Info("User's tasks:", tasks2)
}
