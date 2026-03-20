package parser

import (
	"time"

	"go.uber.org/zap"
)

// TimeParser parses user-provided date/time strings and returns a due time and task text.
type TimeParser struct {
	Logger  *zap.SugaredLogger
	NowFunc func() time.Time
}
