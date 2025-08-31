package store

import (
	"sync"

	"github.com/Maltide/notification_bot_tg/pkg/types"
)

type MemoryStore struct {
	mu     sync.RWMutex
	nextID int64
	data   map[int64]map[int64]types.Task
}
