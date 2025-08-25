package main

import (
	"os"

	"github.com/Maltide/notification_bot_tg/pkg/config"
	"github.com/Maltide/notification_bot_tg/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		os.Exit(1)
	}
	log, err := logger.Logger(cfg.LogLevel)
	if err != nil {
		os.Exit(1)
	}
	log.Debug("Это оно?")
	log.Warn("And what we have here?")
	log.Info("aaand here?")
}
