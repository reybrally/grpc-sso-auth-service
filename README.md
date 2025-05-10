# Сервис авторизации SSO (gRPC + SQLite)

**SSO (Single Sign-On)** — микросервис для аутентификации и выдачи токенов доступа.
Реализован на Go с gRPC-интерфейсом и хранением данных в SQLite.

## 🚀 Возможности

* Регистрация пользователя (`Register`)
* Авторизация и выдача JWT (`Login`)
* Проверка роли администратора (`IsAdmin`)
* Управление схемой через SQL-миграции
* Высокая производительность gRPC/HTTP2

## ⚙️ Установка и запуск

1. Клонировать репозиторий
   git clone [https://github.com/your-username/grpc-sso-auth-service.git](https://github.com/your-username/grpc-sso-auth-service.git)
   cd grpc-sso-auth-service

2. Установить зависимости Go
   go mod download

3. Генерация кода из `.proto` (если нет файлов в `gen/`)
   protoc -I protos protos/\*.proto
   \--go\_out=gen/go --go\_opt=paths=source\_relative
   \--go-grpc\_out=gen/go --go-grpc\_opt=paths=source\_relative

4. Применить миграции
   go run ./cmd/migrator --storage-path=storage/sso.db --migrations-path=migrations
   или
   task migrate

5. Запустить сервис
   go run ./cmd/sso --config=config/local.yaml
   Сервис будет слушать на localhost:44044

## 🔗 RPC-методы

Все методы доступны по gRPC на grpc://localhost:44044

### 1. Register — регистрация пользователя

* Метод: Auth/Register
* Message:
  {
  "email": "[user@example.com](mailto:user@example.com)",
  "password": "secret"
  }
* Response:
  {
  "user\_id": 1
  }

### 2. Login — авторизация и получение токена

* Метод: Auth/Login
* Message:
  {
  "email": "[user@example.com](mailto:user@example.com)",
  "password": "secret",
  "user\_id": 1
  }
* Response:
  {
  "token": "\<JWT\_TOKEN>"
  }

### 3. IsAdmin — проверка прав администратора

* Метод: Auth/IsAdmin
* Message:
  {
  "user\_id": 1
  }
* Response:
  {
  "is\_admin": false
  }

## 📂 Структура проекта

grpc-sso-auth-service/
├ cmd/
│ ├ migrator/      бинарь для миграций
│ └ sso/           точка входа сервиса
├ config/          файлы конфигурации
├ internal/
│ ├ app/           инициализация
│ ├ services/      бизнес-логика
│ └ storage/       репозитории (SQLite)
├ migrations/      SQL-миграции
├ protos/          .proto-файлы
└ gen/             сгенерированный код

## 📝 Конфигурация

config/local.yaml:
env: "local"
storage\_path: "./storage/sso.db"
token\_ttl: 1h0m0s
grpc:
port: 44044
timeout: 10h0m0s

## 🔧 Миграции

* migrations/1\_init.up.sql — создание users и apps
* migrations/2\_add\_is\_admin\_to\_users\_tbl.up.sql — добавление is\_admin
* \*.down.sql для отката

Применение:
task migrate
или
go run ./cmd/migrator --storage-path=storage/sso.db --migrations-path=migrations

