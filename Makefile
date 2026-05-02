.PHONY: help up down build logs clean dev status follow

# Переменные
COMPOSE = docker-compose

help:
	@echo "=========================================="
	@echo "       Service Park - Commands"
	@echo "=========================================="
	@echo ""
	@echo "[ALL SERVICES]"
	@echo "  make up          - Start all services"
	@echo "  make down        - Stop all services"
	@echo "  make build       - Build all services"
	@echo "  make restart     - Restart all services"
	@echo "  make logs        - View all logs"
	@echo "  make ps          - Show container status"
	@echo "  make clean       - Remove all containers and volumes"
	@echo ""
	@echo "[SERVICE PARK ONLY]"
	@echo "  make park-up     - Start service_park only"
	@echo "  make park-down   - Stop service_park only"
	@echo "  make park-build  - Build service_park only"
	@echo "  make park-logs   - View service_park logs"
	@echo "  make park-shell  - Enter service_park database"
	@echo "  make park-restore- Restore database from backup"
	@echo "  make park-backup - Create database backup"
	@echo "  make park-reset  - Reset database (drop and recreate)"
	@echo ""
	@echo "[DEVELOPMENT]"
	@echo "  make dev         - Start services in development mode"
	@echo "  make status      - Show detailed status"
	@echo "  make follow      - Follow all logs"

# Все сервисы
up:
	$(COMPOSE) up -d
	@echo "[OK] All services started"
	@echo "API: http://localhost:8080/api/miinstance"
	@echo "Health: http://localhost:8080/health"

down:
	$(COMPOSE) down
	@echo "[OK] All services stopped"

build:
	$(COMPOSE) build
	@echo "[OK] All services built"

restart: down up
	@echo "[OK] All services restarted"

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

status:
	$(COMPOSE) ps
	@echo ""
	@echo "API Status:"
	@curl -s http://localhost:8080/health 2>/dev/null || echo "API not responding"

follow:
	$(COMPOSE) logs -f --tail=100

clean:
	$(COMPOSE) down -v
	@echo "[OK] Cleaned - all containers and volumes removed"

# Service Park только
park-up:
	$(COMPOSE) up -d service_park_db service_park_app
	@echo "[OK] Service Park started"
	@echo "API: http://localhost:8080/api/miinstance"

park-down:
	$(COMPOSE) stop service_park_db service_park_app
	@echo "[OK] Service Park stopped"

park-build:
	$(COMPOSE) build service_park_app
	@echo "[OK] Service Park built"

park-restart: park-down park-up
	@echo "[OK] Service Park restarted"

park-logs:
	$(COMPOSE) logs -f service_park_app

park-shell:
	@echo "Connecting to database..."
	@echo "Commands: \dt (show tables), SELECT * FROM miinstance LIMIT 10;"
	@echo "Type \q to exit"
	@echo ""
	docker exec -it service_park_db psql -U service_park_user -d service_park_db

park-restore:
	@echo "[RESTORE] Restoring database from backup..."
	@docker cp backup/service_park_backup service_park_db:/tmp/backup.dump 2>/dev/null || echo "Backup file not found"
	@docker exec service_park_db pg_restore -U service_park_user -d service_park_db --no-owner --clean --if-exists /tmp/backup.dump 2>/dev/null || true
	@echo "[OK] Database restored"
	@echo "Verify with: make park-shell"

park-backup:
	@mkdir -p backup
	@docker exec service_park_db pg_dump -U service_park_user -d service_park_db --format=custom > backup/backup_$$(date +%Y%m%d_%H%M%S).dump
	@echo "[OK] Backup created: backup/backup_$$(date +%Y%m%d_%H%M%S).dump"

park-reset:
	@echo "[WARNING] This will DELETE all data in service_park_db"
	@docker exec service_park_db psql -U service_park_user -d postgres -c "DROP DATABASE IF EXISTS service_park_db;"
	@docker exec service_park_db psql -U service_park_user -d postgres -c "CREATE DATABASE service_park_db;"
	@make park-restore
	@echo "[OK] Database reset completed"

dev:
	@echo "Starting development mode..."
	$(COMPOSE) up -d
	@echo ""
	@echo "Follow logs with: make follow"
	@echo "Stop with: make down"