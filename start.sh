#! /bin/sh

set -e

echo "Run DB migrations"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "Start the app"
exec "$@"
