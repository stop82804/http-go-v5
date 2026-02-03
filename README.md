# HTTP Server & Client на Go

HTTP сервер та клієнт для симуляції роботи клієнт-серверної архітектури з логуванням пристроїв мережі.

## Особливості

- **HTTP Сервер** на `localhost:8080`
- **5 HTTP методів**: GET, POST, DELETE, PUT, PATCH
- **Логування пристроїв** з timestamp
- **Кросплатформність**: Windows, Linux, macOS
- **3 способи тестування**: Go клієнт, Bash скрипт, Batch скрипт

## Структура проекту

```
http-go-v5/
├── server.go           # HTTP сервер
├── client.go           # Go клієнт для тестування
├── test_curl.sh        # Bash скрипт для Linux/Mac
├── test_curl.bat       # Batch скрипт для Windows
├── go.mod              # Go модуль
└── README.md           # Документація
```

## Вимоги

- Go 1.21 або новіше
- curl (для скриптів тестування)

## Встановлення та запуск

### 1. Запуск сервера

```bash
go run server.go
```

Сервер запуститься на `http://localhost:8080/`

### 2. Тестування

#### Варіант А: Go клієнт (в окремому терміналі)

```bash
go run client.go
```

#### Варіант Б: CURL скрипт для Linux/macOS

```bash
chmod +x test_curl.sh
./test_curl.sh
```

#### Варіант В: CURL скрипт для Windows

```cmd
test_curl.bat
```

## API Endpoints

### GET /
Повертає вміст лог-файлу з усіма записами про пристрої.

```bash
curl http://localhost:8080/
```

### POST /
Додає новий запис до лог-файлу.

**Параметри:**
- `device_name` - назва пристрою
- `device_type` - тип пристрою (Router, Switch, Server, тощо)
- `ip_address` - IP-адреса пристрою
- `routing_type` - тип маршрутизації (Static, Dynamic, BGP, OSPF, тощо)

**Приклад:**
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "device_name": "Router-Main",
    "device_type": "Router",
    "ip_address": "192.168.1.1",
    "routing_type": "Static"
  }'
```

### DELETE /
Очищає вміст лог-файлу.

```bash
curl -X DELETE http://localhost:8080/
```

### PUT /
Повне оновлення/заміна всього вмісту лог-файлу.

**Приклад:**
```bash
curl -X PUT http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '[
    {
      "device_name": "New-Router",
      "device_type": "Router",
      "ip_address": "10.10.10.1",
      "routing_type": "OSPF"
    }
  ]'
```

### PATCH /
Часткове оновлення - додає нові записи без очищення існуючих.

**Приклад:**
```bash
curl -X PATCH http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '[
    {
      "device_name": "Laptop-Admin",
      "device_type": "Client",
      "ip_address": "192.168.1.150",
      "routing_type": "DHCP"
    }
  ]'
```

## Формат логу

Кожен запис у лог-файлі має формат:
```
[2026-02-03 12:34:56] Пристрій: Router-Main, Тип: Router, IP: 192.168.1.1, Маршрутизація: Static
```

## Тестові сценарії

Скрипти виконують 10 запитів:

1. **DELETE** - очищення лог-файлу
2-5. **POST** - додавання 4 пристроїв (Router, Switch, Server, Firewall)
6. **GET** - перегляд лог-файлу
7. **PATCH** - додавання 2 пристроїв (Laptop, Camera)
8. **GET** - перевірка після PATCH
9. **PUT** - повна заміна на 2 нові пристрої
10. **GET** - фінальна перевірка

## Файли

- `network_devices.log` - автоматично створюється при запуску сервера

## Автор

Sergiy Scherbakov

## Ліцензія

Проект створено для навчальних цілей.