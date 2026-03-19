package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

var moscowLoc *time.Location

func init() {
	var err error
	moscowLoc, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		// In minimal Docker images (e.g. Alpine without tzdata) LoadLocation can fail.
		// Moscow doesn't use DST, so a fixed +03:00 offset is sufficient.
		moscowLoc = time.FixedZone("Europe/Moscow", 3*60*60)
	}
}

// NewPostgresStore constructs a Postgres-backed Store implementation.
func NewPostgresStore(db *sql.DB, log *zap.SugaredLogger) *PostgresStore {
	return &PostgresStore{db: db, log: log}
}

// CreateTask inserts a new task for a user and returns the created task
// with assigned ID and user_task_id. This operation is done in a transaction
// and protected by a mutex to avoid concurrent user_task_id conflicts.
func (ps *PostgresStore) CreateTask(ctx context.Context, t types.Task) (types.Task, error) {
	ps.log.Infof("CreateTask requested: user=%d chat=%d text=%q due=%s", t.UserID, t.ChatID, t.Text, t.DueAt.Format("02.01.06 15:04"))

	ps.mu.Lock()
	defer ps.mu.Unlock()

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return types.Task{}, err
	}
	defer tx.Rollback()

	var nextLocal int64

	err = tx.QueryRowContext(ctx,
		"SELECT COALESCE(MAX(user_task_id),0)+1 FROM tasks WHERE user_id=$1",
		t.UserID,
	).Scan(&nextLocal)
	if err != nil {
		return types.Task{}, err
	}

	t.UserTaskID = nextLocal
	ps.log.Debugf("assigning user_task_id=%d for user=%d", nextLocal, t.UserID)

	err = tx.QueryRowContext(ctx,
		`INSERT INTO tasks (user_id, chat_id, user_task_id, text, due_at)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id`,
		t.UserID, t.ChatID, t.UserTaskID, t.Text, t.DueAt,
	).Scan(&t.ID)
	if err != nil {
		return types.Task{}, err
	}

	if err := tx.Commit(); err != nil {
		return types.Task{}, err
	}

	ps.log.Infof("CreateTask completed: id=%d user=%d user_task_id=%d", t.ID, t.UserID, t.UserTaskID)
	return t, nil
}

// ListTasks returns all tasks for a user ordered by user_task_id.
func (ps *PostgresStore) ListTasks(ctx context.Context, userID int64) ([]types.Task, error) {
	ps.log.Infof("ListTasks requested: user=%d", userID)

	rows, err := ps.db.QueryContext(ctx,
		"SELECT id, user_id, chat_id, user_task_id, text, due_at FROM tasks WHERE user_id=$1 ORDER BY user_task_id",
		userID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]types.Task, 0)

	for rows.Next() {
		var t types.Task
		err := rows.Scan(&t.ID, &t.UserID, &t.ChatID, &t.UserTaskID, &t.Text, &t.DueAt)
		if err != nil {
			return nil, err
		}
		// DB returns timestamptz in UTC; convert to Moscow for display
		t.DueAt = t.DueAt.In(moscowLoc)
		out = append(out, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return []types.Task{}, fmt.Errorf("заметки отсутствуют")
	}

	ps.log.Infof("ListTasks returned %d tasks for user=%d", len(out), userID)
	return out, nil
}

// DeleteTask deletes a user's task by user_task_id and shifts remaining user_task_id
// values down to keep numbering contiguous.
func (ps *PostgresStore) DeleteTask(ctx context.Context, userID, userTaskID int64) error {
	ps.log.Infof("DeleteTask requested: user=%d user_task_id=%d", userID, userTaskID)

	ps.mu.Lock()
	defer ps.mu.Unlock()

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		"DELETE FROM tasks WHERE user_id=$1 AND user_task_id=$2",
		userID, userTaskID,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("task %d not found", userTaskID)
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE tasks SET user_task_id = user_task_id - 1 WHERE user_id=$1 AND user_task_id > $2",
		userID, userTaskID,
	)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	ps.log.Infof("DeleteTask completed: user=%d deleted_user_task_id=%d", userID, userTaskID)

	return nil
}

// NextTask returns the next task to execute (the earliest due_at). If there are no tasks,
// it returns ErrNoTasks.
func (ps *PostgresStore) NextTask(ctx context.Context) (types.Task, error) {
	ps.log.Debugf("NextTask started")

	var t types.Task

	row := ps.db.QueryRowContext(ctx,
		"SELECT id, user_id, chat_id, user_task_id, text, due_at FROM tasks ORDER BY due_at LIMIT 1",
	)
	err := row.Scan(&t.ID, &t.UserID, &t.ChatID, &t.UserTaskID, &t.Text, &t.DueAt)
	if err != nil {
		if err == sql.ErrNoRows {
			ps.log.Debugf("NextTask: no tasks available")
			return types.Task{}, ErrNoTasks
		}
		return types.Task{}, err
	}
	// convert to Moscow for display
	t.DueAt = t.DueAt.In(moscowLoc)
	if t.DueAt.IsZero() {
		return types.Task{}, fmt.Errorf("find task with zero time")
	}
	if t.DueAt.Before(time.Now()) {
		// scheduler умеет обработать задачи "в прошлом" (отправит сразу и удалит)
		ps.log.Infof("NextTask found (in past): id=%d user=%d due=%s", t.ID, t.UserID, t.DueAt.Format("02.01.06 15:04"))
		return t, nil
	}
	ps.log.Infof("NextTask found: id=%d user=%d due=%s", t.ID, t.UserID, t.DueAt.Format("02.01.06 15:04"))
	return t, nil
}
