package store

import (
	"database/sql"
	"sync"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

// MemoryStore is an in-memory Store implementation (intended for tests/dev).
type MemoryStore struct {
	Logger *zap.SugaredLogger
	mu     sync.RWMutex
	nextID int64
	Data   map[int64]map[int64]types.Task
}

// PostgresStore is a Postgres-backed Store implementation.
type PostgresStore struct {
	db  *sql.DB
	log *zap.SugaredLogger
	mu  sync.Mutex
}
