package scheduler

import (
	"context"
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
	Store  types.Store        // источник задач; ожидается, что NextTask вернёт ErrNoTasks при пустом расписании.
	Notify NotifyFunc         // функция доставки уведомлений (Telegram, лог и т.п.).
	Logger *zap.SugaredLogger // общий логгер для отладки.

	refreshCh chan struct{}    // канал сигналов о том, что расписание изменилось.
	timer     *time.Timer      // активный таймер до ближайшей задачи.
	now       func() time.Time // точка расширения для тестов (можно подменить clock).
	// TODO: добавить канал остановки и sync.WaitGroup для graceful shutdown на следующих этапах.
}

// NewTimerScheduler подготавливает структуру и создаёт вспомогательные каналы.
// TODO: добавить параметры конфигурации (буфер канала, дефолтные таймауты) после первых прототипов.
func NewTimerScheduler(store types.Store, notify NotifyFunc, logger *zap.SugaredLogger) *TimerScheduler {
	return &TimerScheduler{
		Store:     store,
		Notify:    notify,
		Logger:    logger,
		refreshCh: make(chan struct{}, 1),
		now:       time.Now,
	}
}

// Start запускает главный цикл обработки напоминаний. Реальную логику студент добавит позже.
func (s *TimerScheduler) Start(ctx context.Context) {
	// TODO: реализовать цикл: получить NextTask, запустить таймер, ждать либо refresh, либо контекст.
}

// Refresh отправляет сигнал в refreshCh, чтобы пересчитать ближайшее напоминание.
func (s *TimerScheduler) Refresh() {
	// TODO: отправить struct{} в refreshCh с защитой от переполнения буфера.
}

// Stop завершает работу планировщика.
func (s *TimerScheduler) Stop(ctx context.Context) error {
	// TODO: корректно остановить таймер и дождаться завершения горутины (graceful shutdown позже).
	return nil
}
