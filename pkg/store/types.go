package store

import (
	"database/sql"
	"sync"

	"go.uber.org/zap"
)

// PostgresStore is a Postgres-backed Store implementation.
type PostgresStore struct {
	db  *sql.DB
	log *zap.SugaredLogger
	mu  sync.Mutex
}
