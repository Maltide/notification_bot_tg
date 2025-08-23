package logger

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var level zapcore.Level

func Logger(loglevel string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	switch loglevel {
	case "debug":
		level = zapcore.DebugLevel

	case "info":
		level = zapcore.InfoLevel

	case "warn":
		level = zapcore.WarnLevel

	case "error":
		level = zapcore.ErrorLevel

	default:
		return nil, errors.New("Непонялничего")
	}

	cfg.Level = zap.NewAtomicLevelAt(level)
	return cfg.Build()
}
