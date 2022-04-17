#!/bin/sh

# stop the script immediately if any command returns non-zero status
set -e

echo "🏃‍♂️ Running database migrations"
/app/migrate -path /app/migration -database "$DB_HOST" -verbose up

echo "🤞 Starting the app..."
# take all parameter passed to the script and run it
exec "$@"