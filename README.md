# URL Shortener

Микросервисное приложение для регистрации пользователей, создания коротких ссылок и сбора аналитики по событиям ссылок.

В проекте три микросервиса:

- **Auth** - регистрация, логин, выдача и проверка JWT.
- **Links** - создание, получение и удаление коротких ссылок. Пишет данные в PostgreSQL, кэширует ссылки в Redis и публикует события в NATS.
- **Analytics** - читает события из NATS, сохраняет их в PostgreSQL и отдает статистику через API.


## Стек

- Go 1.25.5
- PostgreSQL
- Redis
- NATS JetStream
- gRPC + gRPC Gateway
- JWT
- Goose
- Google Wire
- Mockery
- Prometheus
- Grafana

## Структура

```text
.
├── api/                         # proto-файлы auth, links, analytics
├── services/
│   ├── auth/                    # Auth service
│   ├── links/                   # Links service
│   └── analytics/               # Analytics service + NATS consumer + metrics
├── shared/                      # configs, generated proto, jwt, cache, broker, events
├── prometheus/prometheus.yaml   # Prometheus scrape config
├── grafana/provisioning/        # Grafana datasource provisioning
├── docker-compose.yaml
├── go.work
└── makefile
```

## Быстрый Запуск

Требования:

- Docker + Docker Compose
- Go нужен только для локальной разработки и генерации
- `shared/bin/goose` должен существовать для migration jobs. Если его нет, выполнить `make install-goose`.

Запуск всего приложения:

```bash
docker compose up -d --build
```

Compose поднимает:

- `url_shortener_auth`
- `url_shortener_links`
- `url_shortener_analytics`
- PostgreSQL для каждого сервиса
- Redis
- NATS
- Prometheus
- Grafana
- одноразовые migration jobs для auth/links/analytics

## Генерация и инструменты

Установка локальных инструментов:

```bash
make install-all
```

Генерация proto-кода:

```bash
make generate-services
```

Генерация mockery-моков:

```bash
make generate-mocks
```

Полная настройка dev-окружения:

```bash
make app-setup
```

## Миграции

Все миграции:

```bash
make migrations-up
make migrations-down
```

По сервисам:

```bash
make migrations-up-auth
make migrations-up-links
make migrations-up-analytics
```

```bash
make migrations-down-auth
make migrations-down-links
make migrations-down-analytics
```

В Docker Compose миграции запускаются автоматически через одноразовые контейнеры:

- `url_shortener_auth_migrate`
- `url_shortener_links_migrate`
- `url_shortener_analytics_migrate`

## Redis

Redis используется в сервисе **Links** как кэш для коротких ссылок.

Основной источник данных - PostgreSQL, но при получении ссылки сервис сначала проверяет Redis. Если ссылка найдена в кэше, ответ возвращается быстрее и запрос в PostgreSQL не выполняется. Если ссылки в Redis нет, сервис читает ее из PostgreSQL и сохраняет в Redis на время `REDIS_TTL_MINUTES`.

## NATS

Links публикует события:

- `links.created`
- `links.fetched`
- `links.deleted`

Analytics consumer читает `links.>` из stream `LINK_EVENTS`, сохраняет события в таблицу `analytics.link_events` и подтверждает сообщения через `Ack`.

Посмотреть события в stream:

```bash
make nats-events
```

## Prometheus И Grafana

Analytics отдает метрики:

```text
http://localhost:9090/metrics
```

Prometheus:

```text
http://localhost:9091
```

Grafana:

```text
http://localhost:3000
```

Логин по умолчанию:

```text
admin / admin
```

Grafana datasource `Prometheus` подключается автоматически через provisioning:

```text
http://url_shortener_prometheus:9090
```

## API

### Auth

```http
POST /auth/register
POST /auth/login
```

gRPC:

```proto
service AuthService {
  rpc Register(RegisterRequest) returns (google.protobuf.Empty);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc VerifyToken(VerifyTokenRequest) returns (VerifyTokenResponse);
}
```

### Links

REST:

```http
POST   /links
GET    /links
GET    /links/{short_link}
GET    /links/info/{short_link}
DELETE /links/{short_link}
```

gRPC:

```proto
service LinksService {
  rpc CreateLink(CreateLinkRequest) returns (CreateLinkResponse);
  rpc GetOriginalLink(GetOriginalLinkRequest) returns (GetOriginalLinkResponse);
  rpc GetUserLinks(GetUserLinksRequest) returns (GetUserLinksResponse);
  rpc GetLinkInfo(GetLinkInfoRequest) returns (LinkInfo);
  rpc DeleteLink(DeleteLinkRequest) returns (google.protobuf.Empty);
}
```

### Analytics

REST:

```http
POST /analytics/events
GET  /analytics/events
GET  /analytics/links/{short_link}/stats
```

gRPC:

```proto
service AnalyticsService {
  rpc RecordEvent(RecordEventRequest) returns (google.protobuf.Empty);
  rpc GetEvents(GetEventsRequest) returns (GetEventsResponse);
  rpc GetLinkStats(GetLinkStatsRequest) returns (GetLinkStatsResponse);
}
```

## Базы Данных

### Auth

```sql
CREATE TABLE auth.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(100) UNIQUE NOT NULL,
    salt BYTEA NOT NULL,
    salt_password_hash BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Links

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

### Analytics

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

## Тесты

Запуск всех тестов:

```bash
make test
```

Так как проект собран через `go.work`, команда запускает тесты по модулям:

```bash
go test ./services/auth/...
go test ./services/links/...
go test ./services/analytics/...
go test ./shared/...
```

### Покрытие unit-тестами

Основные бизнес-слои покрыты unit-тестами по микросервисам:

- `services/auth` - handler, service и repository слой. Проверяются регистрация, логин, проверка JWT, ошибки валидации, ошибки репозитория и сценарии с невалидными учетными данными.
- `services/links` - handler, service и repository слой. Проверяются создание short link, получение original link, получение информации о ссылке, список ссылок пользователя, удаление ссылки, работа с Redis cache и отправка analytics events через producer.
- `services/analytics` - handler, service, repository и app слой. Проверяются `RecordEvent`, `GetEvents`, `GetLinkStats`, запись событий в БД, получение статистики и обработка событий из NATS.
- `shared` - общие пакеты: JWT generate/verify, gRPC auth interceptor, конфиги, валидация email/password, генерация short link, Redis cache и базовая логика NATS producer/consumer.

В тестах не поднимаются реальные PostgreSQL, Redis и NATS:

- для handler/service слоев используются моки, сгенерированные через `mockery`;
- для repository слоя используется `pgxmock`;
- для Redis cache используется `redismock`;

## Зависимости

- `github.com/jackc/pgx/v4` - PostgreSQL driver.
- `github.com/georgysavva/scany` - scan rows в структуры по `db` тегам.
- `github.com/golang-jwt/jwt/v5` - JWT.
- `golang.org/x/crypto/bcrypt` - password hashing.
- `github.com/go-playground/validator` - validation.
- `github.com/google/wire` - dependency injection.
- `github.com/joho/godotenv` - загрузка `.env`.
- `github.com/stretchr/testify` - assertions и mock helpers.
- `github.com/vektra/mockery/v2` - генерация моков.
- `github.com/pashagolub/pgxmock` - repository tests без настоящей БД.
