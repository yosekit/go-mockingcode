# MockingCode

Backend-as-a-Service platform (аналог Firebase/Supabase) на Go с микросервисной архитектурой.

## 🏗️ Архитектура

```
┌─────────────┐
│   Gateway   │ :8080 - Единая точка входа (API Gateway)
└─────┬───────┘
      │
      ├──────────┐
      │          │
┌─────▼───┐  ┌──▼──────┐  ┌────────┐
│  Auth   │  │ Project │  │  Data  │
│ :8081   │  │  :8082  │  │ :8083  │
└─────┬───┘  └────┬────┘  └────┬───┘
      │           │             │
┌─────▼───────────▼───┐    ┌───▼─────┐
│    PostgreSQL       │    │ MongoDB │
│      :5432          │    │  :27017 │
└─────────────────────┘    └─────────┘
```

### Микросервисы

- **Gateway** (8080) - API Gateway, CORS, аутентификация, проксирование
- **Auth Service** (8081) - регистрация, логин, JWT токены
- **Project Service** (8082) - управление проектами и схемами (коллекциями)
- **Data Service** (8083) - CRUD для моковых данных

### Базы данных

- **PostgreSQL** - хранение пользователей, проектов, схем
- **MongoDB** - хранение моковых данных
- **Redis** - кеширование (готов к использованию)

## 🚀 Быстрый старт

### Предварительные требования

- Docker Desktop
- Windows Terminal с Ubuntu WSL2
- Cursor IDE (опционально)

### Запуск через автоматический скрипт

**Windows:**
```bash
# Запуск окружения (Docker + WSL + Cursor + сервисы)
.\start-dev-environment.bat

# Для остановки введите 'exit' в окне скрипта
```

### Ручной запуск

```bash
# 1. Запустите Docker Desktop

# 2. Перейдите в директорию проекта
cd /usr/projects/mockingcode

# 3. Запустите все сервисы
docker-compose -f docker/docker-compose.dev.yml up -d

# 4. Проверьте статус
docker-compose -f docker/docker-compose.dev.yml ps

# 5. Посмотрите логи
docker-compose -f docker/docker-compose.dev.yml logs -f
```

### Остановка сервисов

```bash
# Остановить все контейнеры
docker-compose -f docker/docker-compose.dev.yml down

# Остановить с удалением volumes (ВНИМАНИЕ: удалит все данные)
docker-compose -f docker/docker-compose.dev.yml down -v
```

## 📡 API Endpoints

### Базовый URL
```
http://localhost:8080
```

### Аутентификация (публичные)

```bash
# Регистрация
POST /auth/register
{
  "email": "user@example.com",
  "password": "password123"
}

# Вход
POST /auth/login
{
  "email": "user@example.com",
  "password": "password123"
}

# Обновление токена
POST /auth/refresh
{
  "refresh_token": "..."
}
```

### Проекты (требуется JWT)

```bash
# Все запросы требуют заголовок:
# Authorization: Bearer {access_token}

# Список проектов
GET /api/projects

# Создать проект
POST /api/projects
{
  "name": "My Project",
  "description": "Project description"
}

# Получить проект
GET /api/projects/{id}

# Обновить проект
PUT /api/projects/{id}
{
  "name": "Updated Name"
}

# Удалить проект
DELETE /api/projects/{id}

# Коллекции проекта
GET /api/projects/{id}/collections
POST /api/projects/{id}/collections
GET /api/projects/{id}/collections/{collectionId}
PUT /api/projects/{id}/collections/{collectionId}
DELETE /api/projects/{id}/collections/{collectionId}
```

### Данные (требуется JWT + API Key)

```bash
# CRUD для документов
GET    /data/{collection}
POST   /data/{collection}
GET    /data/{collection}/{id}
PUT    /data/{collection}/{id}
DELETE /data/{collection}/{id}
```

## 🧪 Тестирование

```bash
# Автоматический тест всех endpoints
./tests/gateway_test.sh

# Ручное тестирование
curl http://localhost:8080/health
```

## 🛠️ Разработка

### Структура проекта

```
mockingcode/
├── auth/           # Сервис аутентификации
├── project/        # Сервис проектов
├── data/           # Сервис данных
├── gateway/        # API Gateway
├── pkg/            # Общие модули
│   └── models/     # Общие модели данных
├── docker/         # Docker конфигурация
│   ├── docker-compose.dev.yml
│   └── .env
├── tests/          # Тесты
└── go.work         # Go workspace
```

### Локальная разработка

```bash
# Запуск отдельного сервиса
cd auth
go run cmd/server/main.go

# Обновление зависимостей
go mod tidy

# Генерация Swagger документации (если нужно)
cd project
swag init -g cmd/server/main.go
```

### Docker Compose

```bash
# Пересобрать конкретный сервис
docker-compose -f docker/docker-compose.dev.yml build gateway

# Перезапустить конкретный сервис
docker-compose -f docker/docker-compose.dev.yml restart gateway

# Посмотреть логи конкретного сервиса
docker-compose -f docker/docker-compose.dev.yml logs -f gateway
```

## 📝 Конфигурация

### Переменные окружения

См. `docker/.env`:

```env
# Базы данных
DB_NAME=mockdb
DB_USER=mockuser
DB_PASSWORD=mockpass
DB_PORT=5432
MONGO_PORT=27017

# Порты сервисов
AUTH_PORT=8081
PROJECT_PORT=8082
DATA_PORT=8083
GATEWAY_PORT=8080

# JWT
JWT_SECRET=your-secret-key
```

## 📚 Документация API

- Gateway: `http://localhost:8080/swagger/`
- Auth: `http://localhost:8081/swagger/`
- Project: `http://localhost:8082/swagger/`
- Data: `http://localhost:8083/swagger/`

## 🎯 План развития

### ✅ Выполнено

- [x] Фаза 1: Настройка окружения (WSL2, Docker, Go workspace)
- [x] Фаза 2.1: Auth Service (регистрация, логин, JWT)
- [x] Фаза 2.2: Project Service (CRUD проектов)
- [x] Фаза 2.4: Management API Gateway
- [x] Фаза 3.1: Project Service - CRUD схем/коллекций
- [x] Фаза 3.2: Data Service - CRUD моковых данных

### 🚧 В процессе

- [ ] Фаза 2.3: gRPC между сервисами
- [ ] Фаза 3.3: Dynamic Router API Gateway (поддомены)
- [ ] Фаза 4: Frontend (Next.js/Nuxt.js)
- [ ] Фаза 5: Генерация данных (Faker), API First, документация

## 🤝 Автоматизация

### Windows скрипты

- `start-dev-environment.bat` - запуск всего окружения
- `start-dev-environment.ps1` - PowerShell версия (с цветным выводом)

Скрипты автоматически:
1. Запускают Docker Desktop
2. Открывают Ubuntu консоль в проекте
3. Запускают Cursor
4. Поднимают все сервисы через docker-compose
5. При вводе `exit` - останавливают всё и очищают

## 📞 Порты

- **8080** - Gateway (единая точка входа)
- **8081** - Auth Service
- **8082** - Project Service
- **8083** - Data Service
- **5432** - PostgreSQL
- **27017** - MongoDB
- **6379** - Redis

## 🔧 Troubleshooting

### Ошибка подключения к БД

```bash
# Проверьте что БД запущены
docker-compose -f docker/docker-compose.dev.yml ps

# Пересоздайте контейнеры БД
docker-compose -f docker/docker-compose.dev.yml up -d postgres mongodb
```

### Порт уже занят

```bash
# Найдите процесс
netstat -ano | findstr :8080

# Остановите процесс или измените порт в .env
```

### Сборка не проходит

```bash
# Очистите старые образы
docker system prune -a

# Пересоберите без кеша
docker-compose -f docker/docker-compose.dev.yml build --no-cache
```

## 📄 Лицензия

MIT

