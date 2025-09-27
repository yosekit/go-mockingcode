#!/bin/bash

BASE_URL="http://localhost:8081"
EMAIL="test2@example.com"
PASS="password123"
WRONG_PASS="_wrongpassword_"

echo "=== Testing Auth Service ==="

# 1. Health check
echo -e "\n1. Testing health endpoint:"
curl -s "$BASE_URL/health" | jq .

# 2. Регистрация нового пользователя
echo -e "\n2. Testing registration:"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'$EMAIL'",
    "password": "'$PASS'"
  }')

echo "$REGISTER_RESPONSE" | jq .

# Извлекаем токены из ответа
ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.access_token')
REFRESH_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.refresh_token')

echo -e "\nAccess Token: $ACCESS_TOKEN"
echo -e "Refresh Token: $REFRESH_TOKEN"

# 3. Попытка регистрации с тем же email (должна быть ошибка)
echo -e "\n3. Testing duplicate registration:"
curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'$EMAIL'",
    "password": "'$PASS'"
  }' | jq .

# 4. Логин с правильными credentials
echo -e "\n4. Testing login with correct credentials:"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'$EMAIL'",
    "password": "'$PASS'"
  }')

echo "$LOGIN_RESPONSE" | jq .

# Обновляем токены
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token')
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.refresh_token')

# 5. Логин с неправильным паролем
echo -e "\n5. Testing login with wrong password:"
curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'$EMAIL'",
    "password": "'$WRONG_PASS'"
  }' | jq .

# 6. Защищенный эндпоинт (должен работать с токеном)
echo -e "\n6. Testing protected endpoint:"
curl -s -X GET "$BASE_URL/protected-route" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# 7. Защищенный эндпоинт без токена (должна быть ошибка)
echo -e "\n7. Testing protected endpoint without token:"
curl -s -X GET "$BASE_URL/protected-route" | jq .

# 8. Обновление токенов
echo -e "\n8. Testing token refresh:"
REFRESH_RESPONSE=$(curl -s -X POST "$BASE_URL/refresh" \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "'$REFRESH_TOKEN'"
  }')

echo "$REFRESH_RESPONSE" | jq .

# 9. Логаут
echo -e "\n9. Testing logout:"
curl -s -X POST "$BASE_URL/logout" \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "'$REFRESH_TOKEN'"
  }' | jq .

# 10. Попытка обновления после логаута (должна быть ошибка)
echo -e "\n10. Testing refresh after logout:"
curl -s -X POST "$BASE_URL/refresh" \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "'$REFRESH_TOKEN'"
  }' | jq .

echo -e "\n=== Test completed ==="