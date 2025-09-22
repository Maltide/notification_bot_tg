package config

import (
	"fmt"
	"os"
)

type Config struct {
	TGToken  string
	LogLevel string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		TGToken:  os.Getenv("TG_TOKEN"),
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
	if cfg.TGToken == "" || cfg.LogLevel == "" {
		return Config{}, fmt.Errorf("empty env")
	}
	return cfg, nil
}
