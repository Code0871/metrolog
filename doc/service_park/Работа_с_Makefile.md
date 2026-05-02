```markdown
# ⚙️ Makefile — управление проектом

В корне проекта `D:\metrolog` находится файл `Makefile`, который содержит команды для удобного управления сервисами. Вместо длинных команд Docker вы используете короткие алиасы.

---

## 🚀 Как использовать

Все команды выполняются из **корня проекта**:

```bash
cd D:\metrolog
make <команда>
```

---

## 📋 Полный список команд

### 📦 Основные команды

| Команда | Что делает | Длинная версия |
|:--------|:-----------|:---------------|
| `make up` | Запускает все сервисы | `docker-compose up -d` |
| `make down` | Останавливает все сервисы | `docker-compose down` |
| `make build` | Пересобирает Docker образы | `docker-compose build` |
| `make restart` | Перезапускает все сервисы | `docker-compose restart` |
| `make logs` | Показывает логи всех сервисов | `docker-compose logs -f` |
| `make ps` | Показывает статус контейнеров | `docker-compose ps` |
| `make clean` | Полностью удаляет контейнеры и тома | `docker-compose down -v` |

### 🐘 Команды для Service Park

| Команда | Что делает |
|:--------|:-----------|
| `make park-up` | Запускает только БД и приложение |
| `make park-down` | Останавливает только Service Park |
| `make park-build` | Пересобирает только приложение |
| `make park-logs` | Показывает логи только приложения |
| `make park-shell` | Входит в консоль PostgreSQL (psql) |
| `make park-restore` | Восстанавливает БД из бэкапа |
| `make park-backup` | Создаёт бэкап текущего состояния БД |
| `make park-reset` | Полностью сбрасывает БД |

### 🛠️ Команды для разработки

| Команда | Что делает |
|:--------|:-----------|
| `make status` | Показывает статус контейнеров и проверяет API |
| `make follow` | Показывает последние 100 строк логов |
| `make dev` | Запускает сервисы в режиме разработки |

---

## 💡 Примеры использования

### ▶️ Первый запуск проекта

```bash
make up
```

### 👀 Посмотреть, что запущено

```bash
make ps
```

### 📜 Посмотреть логи

```bash
make logs            # логи всех сервисов
make park-logs       # только логи приложения
```

### 🐘 Войти в базу данных

```bash
make park-shell
```

**Внутри psql можно выполнять SQL запросы:**

```sql
\dt                              -- показать все таблицы
SELECT * FROM miinstance LIMIT 5; -- показать 5 записей
\q                               -- выйти
```

### 💾 Создать бэкап

```bash
make park-backup
```

Бэкап сохранится в папку `backup/` с именем `backup_20260502_143000.dump`

### 🔄 Восстановить бэкап

```bash
make park-restore
```

### 🔁 Перезапустить только приложение (после изменений кода)

```bash
make park-restart
```

### 🧹 Полностью очистить всё и начать заново

```bash
make clean && make up && make park-restore
```

---

## 🆕 Добавление новых команд

Чтобы добавить новую команду, отредактируйте `Makefile`:

```makefile
my-command:
	docker-compose exec service_park_app sh -c "go test ./..."
	@echo "✅ Tests completed"
```

Затем выполните:

```bash
make my-command
```

---

## 📊 Преимущества использования Makefile

| Без Makefile | С Makefile |
|:-------------|:-----------|
| `docker-compose up -d` | `make up` |
| `docker-compose logs -f service_park_app` | `make park-logs` |
| `docker exec -it service_park_db psql -U service_park_user -d service_park_db` | `make park-shell` |
| `docker cp backup/... && docker exec pg_restore ...` | `make park-restore` |
| ❌ Нужно помнить длинные команды | ✅ Короткие и понятные |
| ❌ Легко ошибиться | ✅ Всегда правильные команды |

---

## ⚡ Быстрая шпаргалка

```bash
make up              # запустить всё
make down            # остановить всё
make park-shell      # войти в БД
make park-logs       # посмотреть логи
make park-restore    # восстановить бэкап
make park-backup     # создать бэкап
make clean           # удалить всё
make help            # показать все команды
```

---

## 📁 Структура Makefile

```makefile
.PHONY: help up down build logs clean

COMPOSE = docker-compose

# Основные команды
up:      $(COMPOSE) up -d
down:    $(COMPOSE) down
build:   $(COMPOSE) build
logs:    $(COMPOSE) logs -f
clean:   $(COMPOSE) down -v

# Команды для Service Park
park-shell:   docker exec -it service_park_db psql -U service_park_user -d service_park_db
park-restore: docker cp backup/service_park_backup service_park_db:/tmp/backup.dump && \
              docker exec service_park_db pg_restore -U service_park_user -d service_park_db --no-owner /tmp/backup.dump
park-backup:  docker exec service_park_db pg_dump -U service_park_user -d service_park_db --format=custom > backup/backup_$$(date +%Y%m%d_%H%M%S).dump
```
```
