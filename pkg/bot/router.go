package bot

import "context"

// Router распределяет входящие сообщения и команды по обработчикам.
// TODO: добавить middleware (логирование, метрики) после первых итераций.
type Router interface {
	HandleMessage(ctx context.Context, upd Update) error
	HandleCommand(ctx context.Context, upd Update) error
}

// SimpleRouter — заготовка базовой реализации.
// TODO: наполнить полями (парсер, store, formatter) и логикой.
type SimpleRouter struct {
	// Parser types.Parser
	// Store  types.Store
	// Logger *zap.SugaredLogger
}

// HandleMessage отвечает за обработку обычных текстовых сообщений (создание напоминаний).
func (r *SimpleRouter) HandleMessage(ctx context.Context, upd Update) error {
	// TODO: вызвать парсер, сохранить задачу, подготовить ответ пользователю.
	return nil
}

// HandleCommand обрабатывает Telegram-команды (/list, /delete, /help).
func (r *SimpleRouter) HandleCommand(ctx context.Context, upd Update) error {
	// TODO: распарсить аргументы, обратиться к стору и подготовить ответ.
	return nil
}
