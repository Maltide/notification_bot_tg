package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

// NotifyFunc delivers a task notification to the end user (e.g. via Telegram).
type NotifyFunc func(ctx context.Context, task types.Task) error

// Scheduler defines the contract for reminder schedulers.
type Scheduler interface {
	Start(ctx context.Context)
	// Refresh просит пересчитать ближайшее напоминание после изменений в стораже.
	Refresh()
	Stop() error
}

// TimerScheduler is a Scheduler implementation based on a single timer.
type TimerScheduler struct {
	Store        types.Store // источник задач;
	Notify       NotifyFunc  // функция доставки уведомлений (Telegram, лог и т.п.).
	Logger       *zap.SugaredLogger
	refreshCh    chan struct{}    // канал сигналов о том, что расписание изменилось.
	timer        *time.Timer      // активный таймер до ближайшей задачи.
	now          func() time.Time // точка расширения для тестов (можно подменить clock).
	current_task types.Task
	wg           sync.WaitGroup
	// ctx          context.Context
	// cancel       context.CancelFunc
	// TODO: добавить канал остановки и sync.WaitGroup для graceful shutdown
}

// NewTimerScheduler constructs a TimerScheduler.
func NewTimerScheduler(store types.Store, notify NotifyFunc, logger *zap.SugaredLogger) *TimerScheduler {
	return &TimerScheduler{
		Store:        store,
		Notify:       notify,
		Logger:       logger,
		refreshCh:    make(chan struct{}, 1),
		now:          time.Now,
		current_task: types.Task{},
		// ctx:          ctx,
		// cancel:       cancel,
	}
}

// Start runs the scheduler main loop. It blocks until ctx is cancelled.
func (s *TimerScheduler) Start(ctx context.Context) {
	if s.timer == nil { // this fake-timer need because of initialization(or we get panicked)
		s.timer = time.NewTimer(time.Hour * 24 * 365)
	}

	s.wg.Add(1)
	defer s.wg.Done()

	for {
		select {

		case <-s.timer.C:
			if s.current_task.UserID == 0 {
				continue
			}
			// s.timer.Stop()
			s.Notify(ctx, s.current_task)

			s.Store.DeleteTask(ctx, s.current_task.UserID, s.current_task.UserTaskID)

			s.Refresh() // give signal to refreshCh => update timer to new task if it exists

			continue

		case <-s.refreshCh:

			if !s.timer.Stop() {
				select {
				case <-s.timer.C:
				default:
				}
			}

			newTask, err := s.Store.NextTask(ctx) // searching for near task

			if err != nil {
				s.Logger.Warn("NextTask error:", zap.Error(err))
				s.current_task = types.Task{}
				s.timer = time.NewTimer(time.Hour * 24 * 365)
				continue
			}

			if newTask.DueAt.Before(time.Now()) {
				s.Notify(ctx, newTask)
				s.Store.DeleteTask(ctx, newTask.UserID, newTask.UserTaskID)
				s.Refresh()
				continue
			}

			s.current_task = newTask // update current_task because this var need for notify users

			s.timer = time.NewTimer(time.Until(newTask.DueAt))

		case <-ctx.Done():

			s.Logger.Info("context done case")

			if !s.timer.Stop() {
				select {
				case <-s.timer.C:
				default:
				}
				s.Logger.Debug("scheduler stopped")
			}

			s.Logger.Info("gorutine was done")

			return
		}
	}
}

// Refresh signals the scheduler to recalculate the nearest reminder.
func (s *TimerScheduler) Refresh() { //use this after calling add,del,upd functions
	select {
	case s.refreshCh <- struct{}{}:
		// сигнал успешно отправлен — планировщик обновится
	default:
		s.Logger.Debug("skip refresh: channel already full")
	}
	// TODO: отправить struct{} в refreshCh с защитой от переполнения буфера.
}

// Stop stops the scheduler.
func (s *TimerScheduler) Stop() error {
	s.timer.Stop()
	return nil
}
