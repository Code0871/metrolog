#!/bin/bash
# entrypoint.sh
yoyo apply --batch --database "$DATABASE_URL" ./migrations
exec "$@"