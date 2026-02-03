#!/bin/bash

echo "========================================"
echo "Тестування HTTP сервера через CURL"
echo "========================================"
echo ""

SERVER="http://localhost:8080/"

# 1. DELETE - Очищення лог-файлу
echo "1. DELETE запит - Очищення лог-файлу"
curl -X DELETE "$SERVER" -w "\nHTTP Status: %{http_code}\n" 2>/dev/null
echo ""
sleep 1

# 2. POST - Додавання Router
echo "2. POST запит - Додавання Router"
curl -X POST "$SERVER" \
  -H "Content-Type: application/json" \
  -d '{
    "device_name": "Router-Main",
    "device_type": "Router",
    "ip_address": "192.168.1.1",
    "routing_type": "Static"
  }' \
  -w "\nHTTP Status: %{http_code}\n" 2>/dev/null
echo ""
sleep 1

# 3. POST - Додавання Switch
echo "3. POST запит - Додавання Switch"
curl -X POST "$SERVER" \
  -H "Content-Type: application/json" \
  -d '{
    "device_name": "Switch-Core-01",
    "device_type": "Switch",
    "ip_address": "192.168.1.2",
    "routing_type": "Dynamic"
  }' \
  -w "\nHTTP Status: %{http_code}\n" 2>/dev/null
echo ""
sleep 1

# 4. POST - Додавання Server
echo "4. POST запит - Додавання Server"
curl -X POST "$SERVER" \
  -H "Content-Type: application/json" \
  -d '{
    "device_name": "Server-DB",
    "device_type": "Server",
    "ip_address": "192.168.1.100",
    "routing_type": "Static"
  }' \
  -w "\nHTTP Status: %{http_code}\n" 2>/dev/null
echo ""
sleep 1

# 5. POST - Додавання Firewall
echo "5. POST запит - Додавання Firewall"
curl -X POST "$SERVER" \
  -H "Content-Type: application/json" \
  -d '{
    "device_name": "Firewall-Edge",
    "device_type": "Firewall",
    "ip_address": "10.0.0.1",
    "routing_type": "BGP"
  }' \
  -w "\nHTTP Status: %{http_code}\n" 2>/dev/null
echo ""
sleep 1

# 6. GET - Читання лог-файлу
echo "6. GET запит - Отримання вмісту лог-файлу"
curl -X GET "$SERVER" -w "\nHTTP Status: %{http_code}\n" 2>/dev/null
echo ""
sleep 1

# 7. PATCH - Часткове оновлення
echo "7. PATCH запит - Часткове оновлення (додавання пристроїв)"
curl -X PATCH "$SERVER" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "device_name": "Laptop-Admin",
      "device_type": "Client",
      "ip_address": "192.168.1.150",
      "routing_type": "DHCP"
    },
    {
      "device_name": "Camera-Entrance",
      "device_type": "IP Camera",
      "ip_address": "192.168.3.20",
      "routing_type": "Static"
    }
  ]' \
  -w "\nHTTP Status: %{http_code}\n" 2>/dev/null
echo ""
sleep 1

# 8. GET - Перевірка після PATCH
echo "8. GET запит - Перевірка після PATCH"
curl -X GET "$SERVER" -w "\nHTTP Status: %{http_code}\n" 2>/dev/null
echo ""
sleep 1

# 9. PUT - Повне оновлення
echo "9. PUT запит - Повне оновлення лог-файлу"
curl -X PUT "$SERVER" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "device_name": "New-Router",
      "device_type": "Router",
      "ip_address": "10.10.10.1",
      "routing_type": "OSPF"
    },
    {
      "device_name": "New-Switch",
      "device_type": "Switch",
      "ip_address": "10.10.10.2",
      "routing_type": "Static"
    }
  ]' \
  -w "\nHTTP Status: %{http_code}\n" 2>/dev/null
echo ""
sleep 1

# 10. GET - Фінальна перевірка
echo "10. GET запит - Фінальна перевірка"
curl -X GET "$SERVER" -w "\nHTTP Status: %{http_code}\n" 2>/dev/null
echo ""

echo "========================================"
echo "Тестування завершено!"
echo "========================================"
