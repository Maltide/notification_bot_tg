package parser

import (
	"time"

	"go.uber.org/zap"
)

type TimeParser struct {
	Logger  *zap.SugaredLogger
	NowFunc func() time.Time
}
