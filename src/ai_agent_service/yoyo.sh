#!/bin/bash
# yoyo.sh - обёртка для yoyo с поддержкой .env

# загружаем config.env
config_file="config/config.env"

if [ -f "$config_file" ]; then
    export $(cat "$config_file" | grep -v '^#' | xargs)
else
    echo "ошибка: $config_file не найден"
    exit 1
fi

# проверяем переменную
db_url="${database_url:-$database_url}"

if [ -z "$db_url" ]; then
    echo "ошибка: database_url не найден в $config_file"
    exit 1
fi

# определяем команду
command="$1"
shift  # убираем команду из аргументов

# для команды new не нужно указывать ./migrations
if [ "$command" = "new" ]; then
    yoyo new "$@" --database "$db_url"
else
    yoyo "$command" "$@" --database "$db_url" ./migrations
fi