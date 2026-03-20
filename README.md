# notification_bot_tg

Telegram-бот для напоминаний: добавляете заметку с датой/временем — бот сохранит её в Postgres и пришлёт сообщение в по записи.

## Возможности

- `/add DD.MM.YY hh:mm Текст` — добавить напоминание
- `/list` — список напоминаний пользователя (локальные ID)
- `/delete <id>` — удалить напоминание по локальному ID
- Хранение задач в Postgres
- Планировщик на одном `time.Timer`: `NextTask → Timer → Notify → Delete → Refresh`

## Быстрый старт (Docker Compose)

1) Создайте `.env` рядом с `docker-compose.yml`:

```dotenv
TG_TOKEN=...
TG_TIMEOUT=10
LOG_LEVEL=debug

DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tasks
DB_HOST=postgres
DB_PORT=5432
```

2) Поднимите сервисы:

```bash
docker compose up --build
```

> Если видите предупреждение про orphan containers, можно чистить старые контейнеры командой:
>
> ```bash
> docker compose down --remove-orphans
> ```

## Локальный запуск (без Docker)

Нужно: Go (версия как в `go.mod`) и доступный Postgres.

1) Экспортируйте переменные окружения (или используйте `.env`):

```bash
export TG_TOKEN=...
export TG_TIMEOUT=10
export LOG_LEVEL=debug

export DB_USER=...
export DB_PASSWORD=...
export DB_NAME=...
export DB_HOST=localhost
export DB_PORT=5432
```

2) Запуск:

```bash
go run ./
```

## Тесты

```bash
go test ./...
```

## Примечания по времени

- Ввод времени в `/add` трактуется как **Europe/Moscow**.
- В минимальных Docker-образах (Alpine без `tzdata`) используется безопасный fallback на фиксированную зону `+03:00`, чтобы не было сдвигов на +3 часа.

## Структура репозитория (кратко)

- `main.go` — сборка зависимостей и запуск бота/шедулера
- `pkg/bot` — Telegram long polling и отправка сообщений
- `pkg/helpers` — обработчики команд (`/add`, `/list`, `/delete`, `/help`)
- `pkg/parser` — парсинг времени и текста для `/add`
- `pkg/store` — Postgres/Memory store и операции над задачами
- `pkg/scheduler` — планировщик напоминаний

Архитектура с диаграммой: см. `ARCHITECTURE.md`.
