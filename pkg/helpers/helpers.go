package helper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Maltide/notification_bot_tg/pkg/messages"
	"github.com/Maltide/notification_bot_tg/pkg/scheduler"
	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

type Handler struct {
	store     types.Store
	scheduler *scheduler.TimerScheduler
	log       *zap.SugaredLogger
	parser    types.Parser
}

func NewHandler(store types.Store, scheduler *scheduler.TimerScheduler, log *zap.SugaredLogger, parser types.Parser) *Handler {
	return &Handler{
		store:     store,
		scheduler: scheduler,
		log:       log,
		parser:    parser,
	}
}

func (h *Handler) HandleCommand(ctx context.Context, userID, chatID int64, command string, args []string) string {
	switch command {
	case "/add":
		return h.handleAdd(ctx, userID, chatID, args)
	case "/list":
		return h.handleList(ctx, userID)
	case "/delete":
		return h.handleDelete(ctx, userID, args)
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

	due, taskText, err := h.parser.ParseAddTask(args)
	if err != nil {
		return fmt.Sprintf("Failed to parse task: %v", err)
	}

	task := types.Task{
		UserID: userID,
		ChatID: chatID,
		Text:   taskText,
		DueAt:  due,
	}

	createdTask, err := h.store.CreateTask(ctx, task)
	if err != nil {
		h.log.Errorf("Failed to create task: %v", err)
		return "Failed to create task"
	}

	h.scheduler.Refresh()
	return messages.MsgAdd(createdTask.ID, due)
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

	return messages.MsgList(tasks)
}

func (h *Handler) handleDelete(ctx context.Context, userID int64, args []string) string {
	if len(args) < 1 {
		return "Usage: /delete <task_id>"
	}

	taskID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return "Invalid task ID"
	}

	if err := h.store.DeleteTask(ctx, userID, taskID); err != nil {
		h.log.Errorf("Failed to delete task: %v", err)
		return "Failed to delete task"
	}

	h.scheduler.Refresh()
	return messages.MsgDelete()
}

func (h *Handler) handleHelp() string {
	return `Available commands:
/add <duration> <text> - Add a new task (e.g., /add 1h30m Buy groceries)
/list - List all your tasks
/delete <id> - Delete a task by ID
/help - Show this help message

Duration format examples: 30s, 5m, 1h, 1h30m, 2h15m30s`
}
