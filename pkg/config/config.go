package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	TGToken   string
	LogLevel  string
	TGTimeout int
}

func LoadConfig() (Config, error) {
	token := os.Getenv("TG_TOKEN")

	loglevel := os.Getenv("LOG_LEVEL")

	timeoutstring := os.Getenv("TG_TIMEOUT")

	if token == "" || loglevel == "" || timeoutstring == "" {
		return Config{}, fmt.Errorf("empty env")
	}

	timeout, err := strconv.Atoi(timeoutstring)
	if err != nil {
		return Config{}, fmt.Errorf("invalid TG_TIMEOUT: %w", err)
	}

	return Config{
		TGToken:   token,
		LogLevel:  loglevel,
		TGTimeout: timeout,
	}, nil
}
