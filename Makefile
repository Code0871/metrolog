.PHONY: help up down db-up db-down app-build app-run app-restart app-stop park-build park-run park-restart park-stop park-logs logs ps restore shell

help:
	@echo Commands:
	@echo ALL TOGETHER:
	@echo   make up          - start everything DB + services
	@echo   make down        - stop everything
	@echo DATABASE ONLY:
	@echo   make db-up       - start only database
	@echo   make db-down     - stop only database
	@echo ALL SERVICES:
	@echo   make app-build   - build all services	@echo   make app-run     - start all services
	@echo   make app-restart - rebuild and restart all
	@echo   make app-stop    - stop all services
	@echo SINGLE SERVICE:
	@echo   make park-build  - build service_park
	@echo   make park-run    - start service_park
	@echo   make park-restart - rebuild and restart service_park
	@echo   make park-stop   - stop service_park
	@echo   make park-logs   - show logs for service_park
	@echo UTILS:
	@echo   make logs        - show all logs
	@echo   make ps          - show container status
	@echo   make restore     - restore database from backup
	@echo   make shell       - enter database console

up:
	docker-compose up -d
	@echo ALL started

douwn:
	docker-compose down
	@echo ALL stopped


restore:
	@echo Restoring backup...
	@docker cp backup/service_park_backup service_park_db:/tmp/backup.dump 2>nul || echo Backup file not found
	@docker exec service_park_db pg_restore -U service_park_user -d service_park_db /tmp/backup.dump 2>nul || echo Restore completed
	@echo Database restored

shell:
	docker exec -it service_park_db psql -U service_park_user -d service_park_db