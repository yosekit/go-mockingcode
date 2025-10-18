# MockingCode API Gateway

API Gateway для платформы MockingCode - единая точка входа для всех микросервисов.

## Архитектура

Gateway объединяет следующие сервисы:
- **Auth Service** (порт 8081) - аутентификация и авторизация
- **Project Service** (порт 8082) - управление проектами и схемами
- **Data Service** (порт 8083) - CRUD операции с моковыми данными

## Endpoints

### Authentication (публичные)

```
POST /auth/register - регистрация пользователя
POST /auth/login    - вход
POST /auth/refresh  - обновление токена
```

### Projects API (защищенные)

```
GET    /api/projects              - список проектов
POST   /api/projects              - создать проект
GET    /api/projects/{id}         - получить проект
PUT    /api/projects/{id}         - обновить проект
DELETE /api/projects/{id}         - удалить проект
GET    /api/projects/{id}/collections         - список коллекций
POST   /api/projects/{id}/collections         - создать коллекцию
GET    /api/projects/{id}/collections/{colId} - получить коллекцию
PUT    /api/projects/{id}/collections/{colId} - обновить коллекцию
DELETE /api/projects/{id}/collections/{colId} - удалить коллекцию
```

### Data API (защищенные)

```
GET    /data/{collection}       - получить документы
POST   /data/{collection}       - создать документ
GET    /data/{collection}/{id}  - получить документ
PUT    /data/{collection}/{id}  - обновить документ
DELETE /data/{collection}/{id}  - удалить документ
```

### Health Check

```
GET /health - проверка работоспособности
```

## Запуск

### Локально

```bash
cd gateway
go run cmd/server/main.go
```

### Docker

```bash
docker-compose -f docker/docker-compose.dev.yml up gateway
```

## Конфигурация

Переменные окружения (`.env`):

```env
GATEWAY_PORT=8080
AUTH_PORT=8081
PROJECT_PORT=8082
DATA_PORT=8083
CORS_ALLOWED_ORIGINS=*
RATE_LIMIT_ENABLED=false
RATE_LIMIT_PER_MIN=100
```

## Аутентификация

Для защищенных endpoints требуется JWT токен в заголовке:

```
Authorization: Bearer {access_token}
```

Публичные endpoints (не требуют токена):
- `/auth/*`
- `/health`
- `/swagger`

## Middleware

- **CORS** - управление cross-origin запросами
- **Auth** - проверка JWT токенов
- **Rate Limiting** (опционально) - ограничение количества запросов

## Пример использования

### Регистрация

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Вход

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Получение проектов

```bash
curl -X GET http://localhost:8080/api/projects \
  -H "Authorization: Bearer {access_token}"
```

## Разработка

### Структура проекта

```
gateway/
├── cmd/
│   └── server/
│       └── main.go           # Точка входа
├── internal/
│   ├── client/               # HTTP клиенты для сервисов
│   │   ├── auth_client.go
│   │   ├── project_client.go
│   │   └── data_client.go
│   ├── config/               # Конфигурация
│   │   └── config.go
│   ├── handler/              # HTTP обработчики
│   │   ├── auth_handler.go
│   │   └── proxy_handler.go
│   ├── middleware/           # Middleware
│   │   ├── auth.go
│   │   └── cors.go
│   └── pkg/
│       └── env/              # Утилиты для env переменных
│           └── env.go
├── Dockerfile
├── go.mod
└── README.md
```

### Добавление нового сервиса

1. Создайте клиент в `internal/client/`
2. Добавьте маршрут в `cmd/server/main.go`
3. Обновите конфигурацию в `internal/config/config.go`
4. Добавьте сервис в `docker-compose.dev.yml`

## TODO

- [ ] Добавить rate limiting
- [ ] Добавить метрики (Prometheus)
- [ ] Добавить трейсинг (Jaeger)
- [ ] Настроить circuit breaker для устойчивости
- [ ] Добавить Swagger документацию

