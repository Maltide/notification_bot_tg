package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/config"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// DBInit opens a Postgres connection, verifies it with a ping, and ensures tables exist.
func DBInit(log *zap.SugaredLogger, cfg config.Config, ctx context.Context) *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	return TableInit(db, log)
}

// TableInit ensures required tables and indexes exist in the database.
// It is idempotent and can be called on every start.
func TableInit(db *sql.DB, log *zap.SugaredLogger) *sql.DB {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		chat_id BIGINT,
		user_task_id BIGINT NOT NULL,
		text TEXT NOT NULL,
		due_at TIMESTAMPTZ NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_tasks_user_user_task_id ON tasks(user_id, user_task_id);
		CREATE INDEX IF NOT EXISTS idx_tasks_due_at ON tasks(due_at);
	`)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
	return db
}
