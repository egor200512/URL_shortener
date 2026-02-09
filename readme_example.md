# PassVault

PassVault — это клиент-серверное решение для безопасного хранения и синхронизации паролей с полным end-to-end шифрованием. Все данные шифруются на стороне клиента перед отправкой на сервер, что гарантирует конфиденциальность даже при компрометации серверной инфраструктуры.

**Сервер** — централизованное хранилище зашифрованных данных с API (gRPC + REST).<br>
**Клиент** — CLI-приложение для шифрования/дешифрования и управления паролями через команды терминала.

## Технологический стек

- **Язык**: Go 1.25.4
- **База данных**: PostgreSQL (physical replication)
- **API**: gRPC + gRPC Gateway (REST)
- **Аутентификация**: JWT токены
- **Контейнеризация**: Docker, Docker Compose
- **Миграции**: Goose
- **Тестирование**: Mockery, testify
- **DI**: Google Wire


## Ключевые особенности:
- End-to-end шифрование
- JWT аутентификация
- PostgreSQL с репликацией для повышения отказоустойчивости
- Полный аудит действий пользователей (с помощью триггеров в PostgreSQL)
- Все слои (handler, service, repo) сервисных endpoint'ов покрыты unit-тестами 


## Структура проекта

```sh
.
├── api/                     # Protocol Buffer декларация
├── cmd/                     # Точки входа в приложение
├── configs/                 # Конфигурационные модули
├── internal/                # Внутренняя бизнес-логика
├── pkg/                     # Переиспользуемые пакеты 
├── migrations/              # SQL миграции 
├── docker-compose.yaml      # Production окружение с репликацией
├── docker-compose_test.yaml # Тестовое окружение с PostgreSQL
├── dockerfile               # Образ приложения
├── entrypoint.sh            # Скрипт запуска контейнера
├── makefile
└── mockery.yaml             # Конфигурация генератора моков
```

## Быстрый старт

### Предварительные требования

- Go 1.25.4+
- Docker и Docker Compose
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
GRPC_HOST="0.0.0.0"
GRPC_PORT="8081" # можно настроить

# HTTP
HTTP_HOST="0.0.0.0"
HTTP_PORT="8080" # можно настроить

# PostgreSQL Master
PG_HOST=postgresql_master
PG_NAME=user 
PG_USER=user 
PG_PASSWORD=pass  # ИЗМЕНИТЕ для production!
PG_PORT=5432 # можно настроить
PG_DSN="host=${PG_HOST} user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_NAME} port=${PG_PORT}"

# PostgreSQL Replica
PG_PORT_REPLICA=5433 # можно настроить
REPL_PASSWORD=repl_pass  # ИЗМЕНИТЕ для production!

# PostgreSQL Tests
PG_PORT_TESTS=5434 # можно настроить

# Migrations
MIGRATION_DIR=./migrations

# JWT
JWT_SECRET_KEY="jwt_secret_key"  # добавьте свой ключ (HS256)
JWT_ACCESS_EXPIRY=60000  # можно настроить (в миллисекундах)
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

Сервер будет доступен на `GRPC_PORT`.

### Запуск клиента (другой терминал)
```sh
cd cmd/client/
go run main.go
```


## Серверная часть

На сервере все записи хранятся в таблицах:

**Данные о пользователях:**
```sql
CREATE TABLE vault.users (
    id UUID PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    salt BYTEA NOT NULL,
    salt_password_hash BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

**Данные о паролях:**
```sql
CREATE TABLE vault.encrypted_data (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES vault.users(id),
    site VARCHAR(100) UNIQUE NOT NULL,
    encrypted_data BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NULL
);
```
Данная таблица хостится на `PG_PORT` и реплицируется на `PG_PORT_REPLICA`.

**Данные об изменениях в таблицах:**
```sql
CREATE TABLE vault.data_audit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    operation CHAR(1) NOT NULL,
    site VARCHAR(100) NOT NULL,
    executed_at TIMESTAMP DEFAULT NOW()
);
```
Создается с помощью `password_audit_trigger`.

## Клиентская часть

Клиент — это консольное приложение, которое работает локально на устройстве пользователя и отвечает за шифрование данных.

Состоит из трёх основных слоёв:

- CLI-интерфейс — принимает команды от пользователя через терминал. Работает в интерактивном режиме.

- gRPC клиент — общается с сервером по gRPC протоколу. Он отправляет зашифрованные данные на сервер и получает их обратно.

- Криптография — здесь происходит всё шифрование и дешифрование данных.

### Как происходит шифрование

Когда зарегистрированный пользователь логинится в систему, генерируются `access_token` и `master_key`, которые хранятся во временном файле `MY_PASS_TMP` локально на устройстве пользователя  и удаляются, как только он выходит из системы.

При сохранении нового пароля, клиент сначала шифрует его локально на устройстве с помощью AES-256-GCM. 

`func EncryptVaultData(data *models.EntryData, masterKey []byte) ([]byte, error)`

```go
type EntryData struct {
	Login    string
	Password string
}
```

Данные отправляются на сервер и хранятся там только в **зашифрованном** виде. Сервер НЕ может прочитать ваш пароль, потому что у него нет ключа для расшифровки(`master_key`).

### Команды для CLI-клиента
Чтобы получить полный список команд и примеры использования, воспользуйтесь командой:
```sh
--help
```
## Взаимодействие клиента и сервера

Общение происходит по GRPC. 

```proto
// Auth
service AuthService {
  // Регистрация нового пользователя в системе
  // Принимает email и пароль, создаёт учётную запись
  // Генерирует соль, хеширует пароль с помощью Argon2
  // Сохраняет данные в таблицу vault.users
  rpc Register(RegisterRequest) returns (google.protobuf.Empty) {...}
  
  // Вход в систему
  // Проверяет email и пароль пользователя
  // При успешной аутентификации возвращает JWT токен
  // Токен используется для авторизации всех последующих запросов
  rpc Login(LoginRequest) returns (LoginResponse) {...}
}

// Vault
service VaultService {
  // Добавление нового пароля
  // Клиент отправляет название сайта и зашифрованные данные (логин + пароль)
  // Сервер сохраняет зашифрованные данные в таблицу vault.encrypted_data
  // Требует JWT токен для авторизации
  rpc AddPassword(AddRequest) returns (google.protobuf.Empty) {...}

  // Получение пароля по названию сайта
  // Клиент запрашивает данные для конкретного сайта
  // Сервер возвращает зашифрованные данные
  // Клиент расшифровывает их локально с помощью мастер-ключа
  rpc GetPassword(GetRequest) returns (GetResponse) {...}

  // Обновление существующего пароля
  // Клиент отправляет новые зашифрованные данные для существующей записи
  // Сервер обновляет поле encrypted_data и updated_at в базе
  // Можно изменить логин, пароль или оба поля
  rpc UpdatePassword(UpdateRequest) returns (google.protobuf.Empty){...}

  // Удаление пароля
  // Удаляет запись из базы данных по названию сайта
  // Действие необратимо, данные восстановить нельзя
  rpc DeletePassword(DeleteRequest) returns (google.protobuf.Empty){...}

// Удаление собственного аккаунта
// Полностью удаляет аккаунт пользователя и все связанные с ним данные
// Включает:
//   - Удаление записи пользователя из vault.users
//   - Каскадное удаление всех паролей из vault.encrypted_data
//   - Удаление записей аудита из vault.data_audit
// Требует действительный JWT токен для авторизации
// Действие необратимо - все данные будут потеряны навсегда
rpc DeleteAccount(DeleteAccountRequest) returns (google.protobuf.Empty){...}
```
После выполнения любого endpoint из `VaultService` происходит добавление записи об этом в `audit_log`.

## Зависимости

Основные библиотеки:
- `google.golang.org/grpc` - gRPC фреймворк
- `github.com/grpc-ecosystem/grpc-gateway/v2` - REST gateway
- `github.com/jackc/pgx/v4` - PostgreSQL драйвер
- `github.com/dgrijalva/jwt-go` - JWT аутентификация
- `github.com/spf13/cobra` - CLI команды
- `github.com/stretchr/testify` - тестирование
- `golang.org/x/crypto` - криптография