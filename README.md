# URL Shortener

URL Shortener — микросервисное приложение для создания и отслеживания коротких ссылок с готовой наблюдаемостью. Логика разнесена на три микросервиса:

- **Auth** — выдаёт JWT и управляет пользователями.  
- **Links** — создаёт/хранит/отдаёт короткие ссылки, кэширует в Redis, шлёт события в NATS.  
- **Analytics** — слушает NATS, пишет события в свою БД и отдаёт агрегаты/метрики.

Из коробки настроены Prometheus и Grafana, есть gRPC и REST, миграции и DI через Wire.

## Технологический стек
- **Язык**: Go 1.25.x  
- **БД**: PostgreSQL (отдельно для auth, links, analytics)  
- **Кэш/Очередь**: Redis, NATS (JetStream)  
- **API**: gRPC + gRPC Gateway (REST)  
- **Аутентификация**: JWT  
- **Наблюдаемость**: Prometheus, Grafana, Postgres Exporter  
- **Миграции**: Goose  
- **DI**: Google Wire  
- **Тесты/моки**: testify, mockery  

## Ключевые особенности
- Короткие ссылки с TTL, хранение в Postgres, кэширование в Redis.
- Auth‑сервис с JWT.
- Analytics‑сервис собирает события (created/fetched/deleted) через NATS и отдаёт агрегированные данные.
- Метрики Prometheus + готовый Grafana дашборд; Postgres‑exporter для аналитической БД.
- gRPC + REST (через gRPC‑Gateway) для всех сервисов.

## Структура проекта
```sh
.
├── api/                     # proto: auth, links, analytics
├── services/
│   ├── auth/                # Auth сервис
│   ├── links/               # Link сервис
│   └── analytics/           # Analytics сервис
├── shared/                  # общие пакеты (broker, configs, models, mocks)
├── deploy/
│   ├── prometheus.yml
│   ├── grafana_user/        # datasources + dashboards (analytics-overview.json)
│   └── postgres_exporter_user_rw/
│       └── queries_analytics.yaml   # кастомные метрики для analytics DB
├── docker-compose.yaml      # dev окружение
├── docker-compose_test.yaml # тестовое окружение
└── makefile
```

## Быстрый старт
### Предварительные требования
https://github.com/egor200512/URL_shortener.git
- Go 1.25.6+
- Docker + Docker Compose
- Make

### Установка

1. Клонируйте репозиторий:
```sh
git clone https://github.com/egor200512/PassVault.git
```

2. Перейдите в папку проекта:
```sh
cd PassVault
```

3. Создайте `.env` файл с переменными окружения:
```sh
# GRPC
GRPC_HOST=localhost
GRPC_AUTH_PORT=8080 # можно настроить
GRPC_LINKS_PORT=8081 # можно настроить
GRPC_ANALYTICS_PORT=8082 # можно настроить

# HTTP
HTTP_HOST=localhost
HTTP_AUTH_PORT=8083 # можно настроить
HTTP_LINKS_PORT=8084 # можно настроить
HTTP_ANALYTICS_PORT=8085 # можно настроить

# PostgreSQL
PG_HOST=localhost
PG_NAME=user # можно настроить
PG_USER=user # можно настроить 
PG_PASSWORD=pass # можно настроить

PG_AUTH_PORT=5432 # можно настроить
PG_LINKS_PORT=5433 # можно настроить
PG_ANALYTICS_PORT=5434 # можно настроить

PG_AUTH_DSN="host=${PG_HOST} user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_NAME} port=${PG_AUTH_PORT}"
PG_LINKS_DSN="host=${PG_HOST} user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_NAME} port=${PG_LINKS_PORT}"
PG_ANALYTICS_DSN="host=${PG_HOST} user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_NAME} port=${PG_ANALYTICS_PORT}"

# PostgreSQL Tests
PG_PORT_TESTS=5435 # можно настроить

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379 # можно настроить
REDIS_PASSWORD=red_pass # можно настроить
REDIS_DB=0
REDIS_TTL_MINUTES=600 # можно настроить

# NATS
NATS_HOST=localhost
NATS_PORT=4222 # можно настроить
NATS_USER=nats_user # можно настроить
NATS_PASSWORD=nats_pass # можно настроить
NATS_URL="nats://${NATS_USER}:${NATS_PASSWORD}@${NATS_HOST}:${NATS_PORT}"
NATS_PREFIX=links
NATS_CONSUMER_NAME="analytics_consumer"
NATS_CONSUMER_BATCH=100
NATS_CONSUMER_WAIT_MINUTES=5 # можно настроить

# Migrations
AUTH_MIGRATION_DIR=./services/auth/migrations
LINKS_MIGRATION_DIR=./services/links/migrations
ANALYTICS_MIGRATION_DIR=./services/analytics/migrations

# JWT
JWT_SECRET_KEY="jwt_secret_key"  # добавьте свой ключ (HS256)
JWT_ACCESS_EXPIRY_MINUTES=60 # можно настроить

# PROMETHEUS
PROMETHEUS_HOST=localhost
PROMETHEUS_PORT=9090 # можно настроить

# GRAFANA
GRAFANA_HOST=localhost
GRAFANA_PORT=3000 # можно настроить

# METRICS
METRICS_HOST=0.0.0.0   
METRICS_ANALYTICS_PORT=9099 # можно настроить
```
Для первого запуска можно просто скопировать и добавить свой `JWT_SECRET_KEY`. 

4. Выполните полную настройку проекта:
```sh
make app-setup
```
Эта команда установит все зависимости, сгенерирует код и запустит тесты.

### Запуск сервера


```sh
docker compose up --build
```



## Наблюдаемость
- Prometheus скрейпит:
  - сервис `analytics` (`/metrics` на `METRICS_ANALYTICS_PORT`)
  - `analytics_pg_exporter` (метрика `events_per_minute` по типам событий).
- Grafana: дашборд `deploy/grafana_user/dashboards/json/analytics-overview.json` (RPS/latency событий, runtime метрики).

## Тестирование
- `go test ./...` — юнит/интеграционные тесты.
- Моки: `mockery` (конфиги в `mockery_*.yaml`).

## Работа ссылочного сервиса
- Создание короткой ссылки: запись в Postgres, кэш в Redis, событие `created` в NATS.
- Получение оригинала: чтение из Redis (фоллбек Postgres), событие `fetched`.
- Удаление/истечение: событие `deleted`.

## Работа аналитики
- Консьюмит события из NATS, пишет в `analytics.link_events`.
- Репозиторий `InsertEvent`/`GetEvents` отдаёт историю, exporter считает `events_per_minute` по типам событий для Grafana.

## Makefile (основные цели)
- `make app-setup` — deps + генерация + миграции.
- `make test` — `go test ./...`.
- `make migrations-up-<service>` — накатывает миграции выбранного сервиса.
- `make cock` — полная перезагрузка контейнеров (осторожно, чистит volumes).

## Перед релизом
- Прогоните `go test ./...` локально.
- Убедитесь, что в `docker-compose.yaml` не остался неиспользуемый `url_shortener_links_pg_exporter`.
- В Grafana нет панелей на несуществующие метрики.
