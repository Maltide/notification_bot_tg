package main

import (
	"notification/pkg/config"
	"notification/pkg/logger"
	"os"
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
}
