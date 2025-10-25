#!/bin/bash

# MockingCode Development Setup Script

set -e

echo "🚀 MockingCode Development Setup"
echo "================================="

# Проверяем наличие .env файла
if [ ! -f .env ]; then
    echo "📝 Создание .env файла из примера..."
    cp .env.example .env
    echo "✅ .env файл создан. Отредактируйте его при необходимости."
else
    echo "✅ .env файл уже существует."
fi

# Проверяем наличие Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker не установлен. Установите Docker для продолжения."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose не установлен. Установите Docker Compose для продолжения."
    exit 1
fi

echo "🐳 Запуск сервисов через Docker Compose..."

# Загружаем переменные окружения
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Запускаем сервисы
docker-compose -f docker/docker-compose.dev.yml up -d

echo "⏳ Ожидание запуска сервисов..."
sleep 10

# Проверяем статус сервисов
echo "🔍 Проверка статуса сервисов..."
docker-compose -f docker/docker-compose.dev.yml ps

# Проверяем доступность API
echo "🌐 Проверка доступности API..."
services=("gateway:8080" "auth:8081" "project:8082" "data:8083")

for service in "${services[@]}"; do
    IFS=':' read -r name port <<< "$service"
    if curl -s "http://localhost:$port/health" > /dev/null; then
        echo "✅ $name сервис доступен на порту $port"
    else
        echo "❌ $name сервис недоступен на порту $port"
    fi
done

echo ""
echo "🎉 Настройка завершена!"
echo ""
echo "📋 Следующие шаги:"
echo "1. Запустите фронтенд: cd frontend && npm install && npm run dev"
echo "2. Откройте браузер: http://localhost:5173"
echo "3. API документация: http://localhost:8080/swagger"
echo ""
echo "🛠️ Полезные команды:"
echo "- Просмотр логов: docker-compose -f docker/docker-compose.dev.yml logs -f"
echo "- Остановка сервисов: docker-compose -f docker/docker-compose.dev.yml down"
echo "- Перезапуск сервисов: docker-compose -f docker/docker-compose.dev.yml restart"
