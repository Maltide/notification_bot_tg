package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger(loglevel string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	ourline := strings.ToLower(strings.TrimSpace(loglevel))

	lvl, err := zapcore.ParseLevel(ourline)
	if err != nil {
		lvl = zapcore.InfoLevel
	}
	cfg.Level = zap.NewAtomicLevelAt(lvl)

	return cfg.Build()
}
