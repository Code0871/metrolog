# Инструкция по запуску Service Park

## Требования

- Docker Desktop (установлен и запущен)
- Go 1.21+
- Git

---

## 1. Запуск базы данных PostgreSQL

### 1.1 Перейдите в корень проекта

```bash
cd D:\metrolog
```

### 1.2 Запустите контейнер с PostgreSQL

```bash
docker-compose up -d
```

### 1.3 Проверьте, что контейнер запущен

```bash
docker-compose ps
```

Ожидаемый вывод:
```
NAME              IMAGE         STATUS          PORTS
service_park_db   postgres:17   Up 2 minutes    0.0.0.0:5432->5432/tcp
```

---

## 2. Восстановление бэкапа

### 2.1 Скопируйте бэкап в контейнер

```bash
docker cp backup/service_park_backup service_park_db:/tmp/backup.dump
```

### 2.2 Восстановите базу данных из бэкапа

```bash
docker exec service_park_db pg_restore -U service_park_user -d service_park_db --no-owner /tmp/backup.dump
```

### 2.3 Проверьте, что данные восстановились

```bash
docker exec -it service_park_db psql -U service_park_user -d service_park_db -c "SELECT COUNT(*) FROM miinstance;"
```

Ожидаемый вывод:
```
  count
---------
 265336
```

---

## 3. Запуск Go сервиса

### 3.1 Перейдите в папку сервиса

```bash
cd D:\metrolog\src\service_park
```

### 3.2 Установите зависимости

```bash
go mod tidy
```

### 3.3 Запустите сервер

```bash
go run main.go
```

Ожидаемый вывод:
```
✅ Database connected successfully
🚀 Server starting on port 8080
```

---

## 4. Эндпоинты API

### Базовый URL: `http://localhost:8080`

| Метод | Эндпоинт | Описание | Пример |
|-------|----------|----------|--------|
| GET | `/api/miinstance` | Список всех СИ (первые 20) | `curl http://localhost:8080/api/miinstance` |
| GET | `/api/miinstance?limit=10&offset=0` | Список с пагинацией | `curl "http://localhost:8080/api/miinstance?limit=10&offset=20"` |
| GET | `/api/miinstance/passport/?idx=` | Получить по паспорту | `curl api/miinstance/passport/?idx=1,10,12,34` |
| GET | `/health` | Проверка здоровья | `curl http://localhost:8080/health` |
| GET | `/ping` | Ping сервера | `curl http://localhost:8080/ping` |
| GET | `/` | Фронтенд (HTML таблица) | Открыть в браузере |

---

## 5. Основные команды Docker

### Управление контейнером

```bash
# Запустить контейнер
docker-compose up -d

# Остановить контейнер
docker-compose down

# Перезапустить контейнер
docker-compose restart

# Посмотреть логи
docker-compose logs -f

# Посмотреть статус
docker-compose ps
```

### Работа с БД

```bash
# Войти в psql консоль
docker exec -it service_park_db psql -U service_park_user -d service_park_db

# Выполнить SQL запрос
docker exec -it service_park_db psql -U service_park_user -d service_park_db -c "SELECT * FROM miinstance LIMIT 5;"

# Создать бэкап текущего состояния
docker exec service_park_db pg_dump -U service_park_user -d service_park_db --format=custom > backup/backup_$(date +%Y%m%d).dump

# Полностью удалить БД (с данными)
docker-compose down -v
```

---

## 6. Проверка работы API

### 6.1 Health check

```bash
curl http://localhost:8080/health
```

Ответ:
```json
{
    "database": "connected",
    "service": "service_park",
    "status": "ok",
    "timestamp": 1746190000
}
```

### 6.2 Получить список СИ

```bash
curl http://localhost:8080/api/miinstance?limit=5
```

Ответ:
```json
{
    "success": true,
    "total": 265336,
    "data": [
        {
            "passport": "14120",
            "name": "Манометр избыточного давления показывающий",
            "type": "МП3-УУ2",
            "state_condition": "В эксплуатации",
            "tech_condition": "Годен",
            "issue_date": "2001-06-01T00:00:00Z",
            "commissioning_date": null,
            "is_fit": true,
            "mpi": 12
        }
    ]
}
```

### 6.3 Получить запись по паспорту

```bash
curl http://localhost:8080/api/miinstance/14120
```

---

## 7. Возможные ошибки и решения

| Ошибка | Решение |
|--------|---------|
| `port is already allocated` | Смените порт в `.env` или остановите другой контейнер: `docker ps` → `docker stop <container>` |
| `relation "miinstance" does not exist` | Восстановите бэкап (шаг 2) |
| `cannot connect to database` | Убедитесь, что Docker запущен: `docker ps` |
| `missing destination name` | Проверьте структуру модели и SQL запрос |

---

## 8. Структура проекта

```
D:\metrolog\
├── backup/
│   └── service_park_backup      # бэкап БД
├── src/
│   └── service_park/
│       ├── config/              # конфигурация
│       ├── db/                  # подключение к БД
│       ├── models/              # структуры данных
│       ├── repository/          # работа с БД
│       ├── service/             # бизнес-логика
│       ├── handlers/            # HTTP обработчики
│       ├── routes/              # маршруты
│       ├── public/              # фронтенд
│       ├── .env                 # переменные окружения
│       └── main.go              # точка входа
├── docker-compose.yaml
└── Makefile
```

---

## 9. Быстрый старт (одной командой)

```bash
# Запустить БД, восстановить бэкап и запустить сервис
cd D:\metrolog && docker-compose up -d && sleep 5 && docker exec service_park_db pg_restore -U service_park_user -d service_park_db /tmp/backup.dump && cd src/service_park && go run main.go
```

---

## 10. Остановка сервиса

```bash
# Остановить Go сервис (Ctrl+C в терминале)

# Остановить контейнер с БД
cd D:\metrolog
docker-compose down
```
