# MockingCode - Система управления Mock-данными

## 🚀 Быстрый старт

### 1. Клонирование и настройка
```bash
git clone <repository-url>
cd mockingcode
```

### 2. Настройка переменных окружения
```bash
# Скопируйте пример конфигурации
cp .env.example .env

# Отредактируйте .env файл под ваши нужды
nano .env
```

### 3. Запуск через Docker Compose
```bash
# Запуск всех сервисов
docker-compose -f docker/docker-compose.dev.yml up -d

# Проверка статуса
docker-compose -f docker/docker-compose.dev.yml ps
```

### 4. Запуск фронтенда
```bash
cd frontend
npm install
npm run dev
```

## 🔧 Конфигурация

### Переменные окружения

#### Frontend (Vite)
- `VITE_API_URL` - URL API Gateway (по умолчанию: `http://localhost:8080`)

#### Gateway
- `GATEWAY_PORT` - Порт Gateway (по умолчанию: `8080`)
- `AUTH_SERVICE_URL` - URL Auth сервиса
- `PROJECT_SERVICE_URL` - URL Project сервиса  
- `DATA_SERVICE_URL` - URL Data сервиса

#### Микросервисы
- `AUTH_PORT` - Порт Auth сервиса (по умолчанию: `8081`)
- `PROJECT_PORT` - Порт Project сервиса (по умолчанию: `8082`)
- `DATA_PORT` - Порт Data сервиса (по умолчанию: `8083`)

#### Базы данных
- `POSTGRES_HOST` - Хост PostgreSQL (по умолчанию: `localhost`)
- `POSTGRES_PORT` - Порт PostgreSQL (по умолчанию: `5432`)
- `MONGO_HOST` - Хост MongoDB (по умолчанию: `localhost`)
- `MONGO_PORT` - Порт MongoDB (по умолчанию: `27017`)

#### Лимиты
- `MAX_PROJECTS_PER_USER` - Максимум проектов на пользователя (по умолчанию: `10`)
- `MAX_COLLECTIONS_PER_PROJECT` - Максимум коллекций на проект (по умолчанию: `20`)
- `MAX_DOCUMENTS_PER_COLLECTION` - Максимум документов на коллекцию (по умолчанию: `500`)

## 🏗️ Архитектура

```
Frontend (React/Preact)
    ↓ HTTP
Gateway (API Gateway + CORS)
    ↓ HTTP
┌─────────────────┬─────────────────┬─────────────────┐
│   Auth Service  │ Project Service │  Data Service   │
│   (JWT)         │ (PostgreSQL)    │ (MongoDB)      │
└─────────────────┴─────────────────┴─────────────────┘
```

## 🚀 Развертывание в продакшене

### 1. Настройка домена
```bash
# В .env файле
PRODUCTION_DOMAIN=yourdomain.com
PRODUCTION_API_URL=https://api.yourdomain.com
```

### 2. SSL сертификаты
```bash
# В .env файле
SSL_CERT_PATH=/etc/ssl/certs/yourdomain.crt
SSL_KEY_PATH=/etc/ssl/private/yourdomain.key
```

### 3. Nginx конфигурация (опционально)
```nginx
server {
    listen 443 ssl;
    server_name api.yourdomain.com;
    
    ssl_certificate /etc/ssl/certs/yourdomain.crt;
    ssl_certificate_key /etc/ssl/private/yourdomain.key;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 🛠️ Разработка

### Hot-reloading для микросервисов
```bash
# Установка air
go install github.com/air-verse/air@latest

# Запуск с hot-reloading
cd data && air
cd project && air  
cd auth && air
cd gateway && air
```

### Тестирование API
```bash
# Проверка здоровья сервисов
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

## 📚 API Документация

- Gateway: http://localhost:8080/swagger
- Auth: http://localhost:8081/swagger
- Project: http://localhost:8082/swagger
- Data: http://localhost:8083/swagger

## 🐛 Отладка

### Логи сервисов
```bash
# Просмотр логов всех сервисов
docker-compose -f docker/docker-compose.dev.yml logs -f

# Логи конкретного сервиса
docker-compose -f docker/docker-compose.dev.yml logs -f gateway
```

### Проверка подключений
```bash
# Проверка доступности сервисов
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

## 📝 Лицензия

MIT License