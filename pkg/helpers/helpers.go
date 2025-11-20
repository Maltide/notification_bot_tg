package helper

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/scheduler"
	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

type Handler struct {
	store     types.Store
	scheduler *scheduler.TimerScheduler
	log       *zap.SugaredLogger
	parser    *types.Parser
}

func NewHandler(store types.Store, scheduler *scheduler.TimerScheduler, log *zap.SugaredLogger, parser types.Parser) *Handler {
	return &Handler{
		store:     store,
		scheduler: scheduler,
		log:       log,
		parser:    &parser,
	}
}

func (h *Handler) HandleCommand(ctx context.Context, userID, chatID int64, command string, args []string) string {
	switch command {
	case "/add":
		return h.handleAdd(ctx, userID, chatID, args)
	case "/list":
		return h.handleList(ctx, userID)
	case "/delete":
		return h.handleDelete(ctx, args)
	case "/help":
		return h.handleHelp()
	default:
		return "Unknown command. Type /help for available commands."
	}
}

func (h *Handler) handleAdd(ctx context.Context, userID, chatID int64, args []string) string {
	if len(args) < 2 {
		return "Usage: /add <duration> <task text>\nExample: /add 1h30m Buy groceries"
	}

	parser.ParseAdd()

	task := types.Task{
		UserID: userID,
		ChatID: chatID,
		Text:   taskText,
		DueAt:  time.Now().Add(duration),
	}

	createdTask, err := h.store.CreateTask(ctx, task)
	if err != nil {
		h.log.Errorf("Failed to create task: %v", err)
		return "Failed to create task"
	}

	h.scheduler.Refresh()
	return fmt.Sprintf("Task added successfully! ID: %d, Will notify you in %s", createdTask.ID, duration)
}

func (h *Handler) handleList(ctx context.Context, userID int64) string {
	tasks, err := h.store.ListTasks(ctx, userID)
	if err != nil {
		h.log.Errorf("Failed to list tasks: %v", err)
		return "Failed to list tasks"
	}

	if len(tasks) == 0 {
		return "You have no pending tasks"
	}

	var sb strings.Builder
	sb.WriteString("Your tasks:\n")
	for _, task := range tasks {
		timeUntil := time.Until(task.DueAt)
		sb.WriteString(fmt.Sprintf("ID: %d | %s | Due in: %s\n", task.ID, task.Text, timeUntil.Round(time.Second)))
	}
	return sb.String()
}

func (h *Handler) handleDelete(ctx context.Context, args []string) string {
	if len(args) < 1 {
		return "Usage: /delete <task_id>"
	}

	taskID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return "Invalid task ID"
	}

	if err := h.store.DeleteTask(ctx, taskID); err != nil {
		h.log.Errorf("Failed to delete task: %v", err)
		return "Failed to delete task"
	}

	h.scheduler.Refresh()
	return "Task deleted successfully"
}

func (h *Handler) handleHelp() string {
	return `Available commands:
/add <duration> <text> - Add a new task (e.g., /add 1h30m Buy groceries)
/list - List all your tasks
/delete <id> - Delete a task by ID
/help - Show this help message

Duration format examples: 30s, 5m, 1h, 1h30m, 2h15m30s`
}
