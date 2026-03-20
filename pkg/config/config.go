package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	TGToken    string
	LogLevel   string
	TGTimeout  int
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

// LoadConfig loads configuration from environment (and optionally .env) and validates
// required variables such as TG_TIMEOUT and TG_TOKEN.
func LoadConfig() (Config, error) {
	// .env удобен локально, но в Docker/CI переменные обычно приходят из окружения.
	_ = godotenv.Load()

	timeoutStr := os.Getenv("TG_TIMEOUT")
	if timeoutStr == "" {
		return Config{}, fmt.Errorf("config: TG_TIMEOUT is empty")
	}

	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return Config{}, fmt.Errorf("config: error converting TG_TIMEOUT to int: %v", err)
	}

	tgToken := os.Getenv("TG_TOKEN")
	if tgToken == "" {
		return Config{}, fmt.Errorf("config: TG_TOKEN is empty")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	return Config{
		TGToken:    tgToken,
		LogLevel:   os.Getenv("LOG_LEVEL"),
		TGTimeout:  timeout,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     dbHost,
		DBPort:     dbPort,
	}, nil
}
