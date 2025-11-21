package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/Maltide/notification_bot_tg/pkg/types"
	"go.uber.org/zap"
)

// NotifyFunc описывает функцию, которая доставляет напоминание конечному пользователю.
// TODO: расширить контракт (форматирование текста, ретраи), когда появится реальный Telegram-адаптер.
type NotifyFunc func(ctx context.Context, task types.Task) error

// Scheduler задаёт общий контракт планировщиков напоминаний.
type Scheduler interface {
	// Start запускает фоновый цикл и должен вызываться в горутине.
	Start(ctx context.Context)
	// Refresh просит пересчитать ближайшее напоминание после изменений в стораже.
	Refresh()
	// Stop завершает работу планировщика и освобождает ресурсы.
	Stop(ctx context.Context) error
}

// TimerScheduler — базовая реализация Scheduler, использующая один таймер.
type TimerScheduler struct {
	Store        types.Store        // источник задач; ожидается, что NextTask вернёт ErrNoTasks при пустом расписании.
	Notify       NotifyFunc         // функция доставки уведомлений (Telegram, лог и т.п.).
	Logger       *zap.SugaredLogger // общий логгер для отладки.
	refreshCh    chan struct{}      // канал сигналов о том, что расписание изменилось.
	timer        *time.Timer        // активный таймер до ближайшей задачи.
	now          func() time.Time   // точка расширения для тестов (можно подменить clock).
	current_task types.Task
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	// TODO: добавить канал остановки и sync.WaitGroup для graceful shutdown на следующих этапах. ????????7
}

// NewTimerScheduler подготавливает структуру и создаёт вспомогательные каналы.
// TODO: добавить параметры конфигурации (буфер канала, дефолтные таймауты) после первых прототипов.
func NewTimerScheduler(store types.Store, notify NotifyFunc, logger *zap.SugaredLogger) *TimerScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &TimerScheduler{
		Store:        store,
		Notify:       notify,
		Logger:       logger,
		refreshCh:    make(chan struct{}, 1),
		now:          time.Now,
		current_task: types.Task{},
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start запускает главный цикл обработки напоминаний. Реальную логику студент добавит позже.
func (s *TimerScheduler) Start(ctx context.Context) {
	// TODO: реализовать цикл: получить NextTask, запустить таймер, ждать либо refresh, либо контекст.
	if s.timer == nil { // this fake-timer need because of initialization(or we get panicked)
		s.timer = time.NewTimer(time.Hour * 24 * 365)
	}
	s.wg.Add(1)
	defer s.wg.Done()
	for {
		select {
		case <-s.timer.C:
			if s.current_task.UserID != 0 {
				s.Notify(ctx, s.current_task)
				s.timer.Stop()
				s.Store.DeleteTask(ctx, s.current_task.UserID, s.current_task.ID)
				s.Refresh() // give signal to refreshCh => update timer to new task if it exists
			}
		case <-s.refreshCh:
			if !s.timer.Stop() {
				<-s.timer.C
			}
			newTask, err := s.Store.NextTask(ctx) // searching for near task
			if err != nil {
				s.Logger.Warn("NextTask error:", zap.Error(err))
				continue
			}
			if newTask.DueAt.Before(time.Now()) {
				s.Notify(ctx, newTask)
				s.Store.DeleteTask(ctx, newTask.UserID, newTask.ID)
				continue
			}
			s.current_task = newTask // update current_task because this var need for notify users
			s.timer = time.NewTimer(time.Until(newTask.DueAt))
		case <-s.ctx.Done():
			s.Logger.Info("context done case")
			if !s.timer.Stop() {
				s.Logger.Debug("timer already stopped")
			}
			s.Logger.Info("gorutine was done")
			return
		}
	}
}

// Refresh отправляет сигнал в refreshCh, чтобы пересчитать ближайшее напоминание.
func (s *TimerScheduler) Refresh() { //use this after calling add,del,upd functions
	select {
	case s.refreshCh <- struct{}{}:
		// сигнал успешно отправлен — планировщик обновится
	default:
		s.Logger.Debug("skip refresh: channel already full")
	}
	// TODO: отправить struct{} в refreshCh с защитой от переполнения буфера.
}

// Stop завершает работу планировщика.
func (s *TimerScheduler) Stop(ctx context.Context) error {
	// TODO: корректно остановить таймер и дождаться завершения горутины (graceful shutdown позже).
	s.cancel()
	s.wg.Wait() //can be useful in the future?
	return nil
}
