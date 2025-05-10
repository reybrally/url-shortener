
# URL Shortener REST API

**URL Shortener REST API** — лёгкий REST-сервис на Go, который позволяет сокращать длинные URL и перенаправлять пользователей по коротким alias.

## 🚀 Основные возможности

- Сокращение URL (POST /urls/)  
- Редирект по alias (GET /{alias})  
- Поддержка кастомного или автоматически сгенерированного alias  
- BasicAuth на POST /urls/  
- Хранение данных в SQLite (встроенная БД)  
- Логирование через slog (pretty-вывод)  
- Unit- и интеграционные тесты

## ⚙️ Быстрый старт

1. Клонировать репозиторий  
   ```bash
   git clone [https://github.com/reybrally/url-shortener-rest-api.git](https://github.com/reybrally/url-shortener.git)
   cd url-shortener
```

2. Установить зависимости Go

   ```bash
   go mod download
   ```

3. Скопировать и настроить конфиг

   ```bash
   cp config/prod.yaml config/local.yaml
   ```

   В файле `config/local.yaml` задать:

   ```yaml
   env: "local"               # local | dev | prod
   storage_path: "./storage/storage.db"
   http_server:
     address: "localhost:8080"
     timeout: 4s
     idle_timeout: 60s
     user: "myuser"           # для BasicAuth
     password: "mypass"       # для BasicAuth
   ```

4. Применить миграции

   ```bash
   go run ./cmd/url-shortener/migrator/main.go \
     --storage-path=storage/storage.db \
     --migrations-path=migrations
   ```

   Или с Taskfile:

   ```bash
   task migrate
   ```

5. Запустить сервис

   ```bash
   go run ./cmd/url-shortener/main.go
   ```

   Сервис слушает на `http://localhost:8080`.

## 🔗 HTTP-API

### 1. Сокращение URL

```http
POST /urls/ HTTP/1.1
Host: localhost:8080
Authorization: Basic bXl1c2VyOm15cGFzcw==   # Base64(myuser:mypass)
Content-Type: application/json

{
  "url": "https://example.com/very/long/path",
  "alias": "myalias"    # необязательно, если не указан — сгенерируется
}
```

* **201 Created**

  ```json
  { "alias": "myalias" }
  ```
* **400 Bad Request** — неверный JSON или URL
* **409 Conflict** — alias уже занят

### 2. Редирект по alias

```http
GET /myalias HTTP/1.1
Host: localhost:8080
```

* **302 Found** с заголовком `Location: https://example.com/very/long/path`
* **404 Not Found** — если alias не существует

## 📂 Структура проекта

```
url-shortener-rest-api/
├── cmd/
│   └── url-shortener/       # main.go
├── config/                  # prod.yaml, копируйте в local.yaml
├── internal/
│   ├── config/              # загрузка конфига через cleanenv
│   ├── http-server/
│   │   ├── handlers/        # redirect и save
│   │   └── middleware/      # logger, BasicAuth и др.
│   ├── lib/
│   │   ├── logger/          # slog pretty
│   │   └── random/          # генерация alias
│   └── storage/
│       └── sqlite/          # репозиторий SQLite
├── migrations/              # SQL-скрипты up/down
├── storage/                 # файл базы storage.db
├── tests/                   # unit & интеграционные тесты
├── go.mod
└── go.sum
```

## 📝 Конфигурация

В `config/local.yaml`:

```yaml
env: "local"
storage_path: "./storage/storage.db"
http_server:
  address: "localhost:8080"
  timeout: 4s
  idle_timeout: 60s
  user: "myuser"
  password: "mypass"
```

## 🧪 Тесты

```bash
go test ./internal/http-server/handlers/...
go test ./internal/storage/sqlite/...
# или все сразу
go test ./...

```
