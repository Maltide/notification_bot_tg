package store

import (
	"sync"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

type MemoryStore struct {
	Logger *zap.SugaredLogger
	mu     sync.RWMutex
	nextID int64
	Data   map[int64]map[int64]types.Task
}
