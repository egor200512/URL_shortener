# URL Shortener

URL Shortener — production-ready микросервисное приложение для быстрого создания и детальной аналитики коротких ссылок. Логика разнесена на три микросервиса:

- **Auth** — выдаёт JWT и управляет пользователями.  
- **Links** — создаёт/хранит/отдаёт короткие ссылки, кэширует в Redis.  
- **Analytics** — читает данные из своей БД и отдаёт агрегаты.

Из коробки есть gRPC и REST, миграции и DI через Wire.

## Технологический стек
- **Язык**: Go 1.25.6 
- **БД**: PostgreSQL  
- **Кэш**: Redis  
- **API**: gRPC + gRPC Gateway (REST)  
- **Аутентификация**: JWT  
- **Миграции**: Goose  
- **DI**: Google Wire  
- **Тесты/моки**: testify, mockery  

## Ключевые особенности
- Короткие ссылки с TTL, хранение в Postgres, кэширование в Redis.
- Auth‑сервис с JWT.
- Analytics‑сервис отдаёт агрегированные данные из своей БД.
- gRPC + REST (через gRPC‑Gateway) для всех сервисов.
- Все слои (handler, service, repo) сервисных endpoint'ов покрыты unit-тестами

## Структура проекта
```sh
.
├── api/                     # proto: auth, links, analytics
├── services/
│   ├── auth/                # Auth сервис
│   ├── links/               # Link сервис
│   └── analytics/           # Analytics сервис
├── shared/                  # общие пакеты (configs, models, mocks)
├── docker-compose.yaml      # dev окружение
├── docker-compose_test.yaml # тестовое окружение
└── makefile
```

## Быстрый старт
### Предварительные требования

- Go 1.25.6+
- Docker + Docker Compose
- Make

### Установка

1. Клонируйте репозиторий:
```sh
git clone https://github.com/egor200512/URL_shortener.git
```

2. Перейдите в папку проекта:
```sh
cd URL_shortener
```

3. Создайте `.env` файл с переменными окружения:
```sh
# GRPC
GRPC_HOST=0.0.0.0
GRPC_AUTH_PORT=8080 # можно настроить
GRPC_LINKS_PORT=8081 # можно настроить
GRPC_ANALYTICS_PORT=8082 # можно настроить
GRPC_AUTH_DOCKER_HOST=url_shortener_auth

# HTTP
HTTP_HOST=0.0.0.0
HTTP_AUTH_PORT=8083 # можно настроить
HTTP_LINKS_PORT=8084 # можно настроить
HTTP_ANALYTICS_PORT=8085 # можно настроить

# PostgreSQL
PG_AUTH_HOST=url_shortener_auth_pg
PG_LINKS_HOST=url_shortener_links_pg
PG_ANALYTICS_HOST=url_shortener_analytics_pg
PG_TESTS_HOST=0.0.0.0
PG_NAME=user # можно настроить
PG_USER=user # можно настроить
PG_PASSWORD=pass # можно настроить

PG_DOCKER_PORT=5432 

PG_AUTH_DSN="host=${PG_AUTH_HOST} user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_NAME} port=${PG_DOCKER_PORT}"
PG_LINKS_DSN="host=${PG_LINKS_HOST} user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_NAME} port=${PG_DOCKER_PORT}"
PG_ANALYTICS_DSN="host=${PG_ANALYTICS_HOST} user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_NAME} port=${PG_DOCKER_PORT}"

# PostgreSQL Tests
PG_PORT_TESTS=5435 # можно настроить

# Redis
REDIS_HOST=url_shortener_redis
REDIS_PORT=6379 # можно настроить
REDIS_PASSWORD=red_pass
REDIS_DB=0
REDIS_TTL_MINUTES=600

# Migrations
AUTH_MIGRATION_DIR=./services/auth/migrations
LINKS_MIGRATION_DIR=./services/links/migrations
ANALYTICS_MIGRATION_DIR=./services/analytics/migrations

# JWT
JWT_SECRET_KEY="jwt_secret_key"  # добавьте свой ключ (HS256)
JWT_ACCESS_EXPIRY_MINUTES=60

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

## Тесты
Запустить unit‑тесты можно следующей командой: 
```sh
make tests
```

## Серверная часть
### Серверная часть состоит из трёх микросервисов.<br>

- **Auth** хранит пользователей в своей БД, выдаёт/проверяет JWT.
- **Links** пишет короткие ссылки в свою БД и кэширует их в Redis.
- **Analytics** считает агрегаты по своей БД и отдаёт их через gRPC/REST.

Сервисы общаются по gRPC; REST добавлен через gRPC‑Gateway.

### Данные приложения хранятся в трёх независимых таблицах.
```sql
CREATE TABLE auth.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(100) UNIQUE NOT NULL,
    salt BYTEA NOT NULL,
    salt_password_hash BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```
```sql
CREATE TABLE links.short_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    short_link VARCHAR(10) UNIQUE NOT NULL,
    original_link_host TEXT NOT NULL,
    original_link TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```
```sql
CREATE TABLE analytics.link_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type TEXT NOT NULL CHECK (event_type IN ('created', 'fetched', 'deleted')),
    user_id UUID NOT NULL,
    short_link VARCHAR(10) NOT NULL,
    original_link TEXT NOT NULL,
    executed_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Основные ручки (REST/gRPC)

```proto
// Auth
service AuthService {
    // Принимает email/пароль, сохранаяет юзера в auth.users
    rpc Register (RegisterRequest) returns (google.protobuf.Empty);
    // Проверяет креды, возвращает JWT access_token
    rpc Login (LoginRequest) returns (LoginResponse);
    // Проверяет валидность токена (используется сервисами)
    rpc VerifyToken (VerifyTokenRequest) returns (VerifyTokenResponse);
}

// Links
service LinksService {
    // Создаёт короткую ссылку, сохраняет в links.short_links
    rpc CreateLink (CreateLinkRequest) returns (CreateLinkResponse);
    // Отдает оригинал по slug
    rpc GetOriginalLink (GetOriginalLinkRequest) returns (GetOriginalLinkResponse);
    // Получает все ссылки пользователя
    rpc GetUserLinks (GetUserLinksRequest) returns (GetUserLinksResponse);
    // Возвращает метаданные по ссылке
    rpc GetLinkInfo (GetLinkInfoRequest) returns (LinkInfo);
    // Удаляет ссылку
    rpc DeleteLink (DeleteLinkRequest) returns (google.protobuf.Empty);
}

// Analytics
service AnalyticsService {
  // Возвращает события/агрегаты из analytics.link_events
  rpc GetEvents(GetEventsRequest) returns (GetEventsResponse);
  // Healthcheck
  rpc Health(google.protobuf.Empty) returns (google.protobuf.Empty);
}
```

## Зависимости
- `google.golang.org/grpc`, `github.com/grpc-ecosystem/grpc-gateway/v2` — gRPC + REST Gateway.
- `github.com/jackc/pgx/v4`, `github.com/georgysavva/scany` — драйвер и маппинг PostgreSQL.
- `github.com/redis/go-redis/v9` — кэширование коротких ссылок в Links.
- `github.com/google/wire` — DI.
- `github.com/pressly/goose/v3` — миграции.
- `github.com/stretchr/testify`, `github.com/vektra/mockery/v3` — тестирование и генерация моков.
