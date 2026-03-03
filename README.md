# Polytech Schedule TMA (Telegram Mini App)

Высоконагруженный сервис расписания для студентов с интеграцией в Telegram Mini Apps. Проект реализует полный цикл работы с расписанием: от автоматического парсинга внешних ресурсов до доставки асинхронных PUSH-уведомлений об отменах и заменах пар.

## 🚀 Ключевые особенности (Features)

- **Telegram Mini App (TMA):** Современный frontend (React + Vite, TailwindCSS) прямо внутри Telegram.
- **Фоновый парсинг (Cron Worker):** Автоматическая синхронизация с источником расписания, отслеживание изменений в реальном времени.
- **Умные уведомления:** Асинхронная рассылка PUSH-уведомлений пользователям при отмене пар или замене преподавателя.
- **Clean Architecture:** Бэкенд строго разделен на слои (Handlers, UseCases, Repositories, Entities), что обеспечивает легкое масштабирование и поддержку.
- **Мониторинг:** Встроенные метрики (Prometheus) и дашборды (Grafana) для отслеживания бизнес-метрик (DAU) и технических показателей (RPS, Latency).

## 🛠 Технологический стек

**Backend:**
- Go 1.25
- [Fiber v3](https://github.com/gofiber/fiber) (Высокопроизводительный HTTP-фреймворк)
- PostgreSQL + GORM
- Telegram Bot API

**Frontend:**
- React 19 + TypeScript
- TailwindCSS
- Sentry (мониторинг ошибок на клиенте)

**Инфраструктура (DevOps):**
- Docker & Docker Compose
- Nginx (Reverse proxy + SSL/Certbot)
- Prometheus + Grafana
- GitHub Actions (CI/CD)

## 📁 Структура проекта (Clean Architecture)

```text
├── cmd/app/              # Точка входа приложения (main.go, Dependency Injection)
├── internal/
│   ├── app/              # Инициализация Fiber, DB, Config, Logger
│   ├── entity/           # Доменные модели базы данных
│   ├── handler/          # HTTP контроллеры (REST API) и Middleware
│   ├── usecase/          # Бизнес-логика (расписание, юзеры, отзывы)
│   ├── repository/       # Слой работы с БД (Postgres / GORM)
│   ├── worker/           # Фоновые процессы (Cron парсер)
│   └── metrics/          # Кастомные метрики Prometheus
├── pkg/                  # Внешние пакеты (парсер XML, http-клиент политеха)
└── frontend/             # Исходный код Telegram Web App (React)
```

## 📈 Мониторинг

В проекте реализован сбор кастомных метрик:
- `polytech_http_requests_total` — счетчик запросов по эндпоинтам.
- `polytech_http_request_duration_seconds` — гистограмма времени ответа API.
- `polytech_users_total` и `polytech_users_active_today` — бизнес-метрики регистраций и ежедневной активности (DAU).
