# URL Shortener

URL Shortener - микросервисное приложение для создания, хранения и получения коротких ссылок.

Логика разнесена на два сервиса:

- **Auth** - регистрация, логин, выдача и проверка JWT.
- **Links** - создание коротких ссылок, хранение в PostgreSQL и кэширование в Redis.

## Технологический стек

- **Язык**: Go 1.25.5
- **БД**: PostgreSQL
- **Кэш**: Redis
- **API**: gRPC + gRPC Gateway (REST)
- **Аутентификация**: JWT
- **Миграции**: Goose
- **DI**: Google Wire

## Структура проекта

```sh
.
├── api/                     # proto: auth, links
├── services/
│   ├── auth/                # Auth сервис
│   └── links/               # Links сервис
├── shared/                  # общие пакеты, configs, cache, jwt
├── docker-compose.yaml      # dev окружение
└── makefile
```

## Быстрый старт

### Требования

- Go 1.25.5+
- Docker + Docker Compose
- Make

### Установка

```sh
git clone https://github.com/egor200512/URL_shortener.git
cd URL_shortener
```

Создайте `.env` файл:

```sh
# GRPC
GRPC_HOST=0.0.0.0
GRPC_AUTH_PORT=8080
GRPC_LINKS_PORT=8081
GRPC_AUTH_DOCKER_HOST=url_shortener_auth

# HTTP
HTTP_HOST=0.0.0.0
HTTP_AUTH_PORT=8083
HTTP_LINKS_PORT=8084

# PostgreSQL
PG_AUTH_HOST=url_shortener_auth_pg
PG_LINKS_HOST=url_shortener_links_pg
PG_NAME=user
PG_USER=user
PG_PASSWORD=pass
PG_AUTH_PORT=5432
PG_LINKS_PORT=5433
PG_DOCKER_PORT=5432

PG_AUTH_DSN="host=${PG_AUTH_HOST} user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_NAME} port=${PG_DOCKER_PORT}"
PG_LINKS_DSN="host=${PG_LINKS_HOST} user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_NAME} port=${PG_DOCKER_PORT}"

# Redis
REDIS_HOST=url_shortener_redis
REDIS_PORT=6379
REDIS_PASSWORD=red_pass
REDIS_DB=0
REDIS_TTL_MINUTES=600

# Migrations
AUTH_MIGRATION_DIR=./services/auth/migrations
LINKS_MIGRATION_DIR=./services/links/migrations

# JWT
JWT_SECRET_KEY="jwt_secret_key"
JWT_ACCESS_EXPIRY_MINUTES=60

```

Сгенерируйте код и установите локальные инструменты:

```sh
make app-setup
```

Запуск:

```sh
docker compose up --build
```

## Сервисы

- **Auth** хранит пользователей в своей БД, выдаёт и проверяет JWT.
- **Links** пишет короткие ссылки в свою БД и кэширует их в Redis.

Сервисы общаются по gRPC; REST добавлен через gRPC Gateway.

## Основные таблицы

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

## Основные ручки

```proto
service AuthService {
    rpc Register (RegisterRequest) returns (google.protobuf.Empty);
    rpc Login (LoginRequest) returns (LoginResponse);
    rpc VerifyToken (VerifyTokenRequest) returns (VerifyTokenResponse);
}

service LinksService {
    rpc CreateLink (CreateLinkRequest) returns (CreateLinkResponse);
    rpc GetOriginalLink (GetOriginalLinkRequest) returns (GetOriginalLinkResponse);
    rpc GetUserLinks (GetUserLinksRequest) returns (GetUserLinksResponse);
    rpc GetLinkInfo (GetLinkInfoRequest) returns (LinkInfo);
    rpc DeleteLink (DeleteLinkRequest) returns (google.protobuf.Empty);
}
```

## Зависимости

- `google.golang.org/grpc`, `github.com/grpc-ecosystem/grpc-gateway/v2` - gRPC + REST Gateway.
- `github.com/jackc/pgx/v4`, `github.com/georgysavva/scany` - PostgreSQL.
- `github.com/redis/go-redis/v9` - кэширование коротких ссылок.
- `github.com/google/wire` - DI.
- `github.com/pressly/goose/v3` - миграции.
