#!/bin/sh
# Wait for postgres to be ready, then run migrations

set -e

echo "Waiting for postgres..."

until PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - running migrations"

# Change to app directory
cd /app

# Check if migrations directory exists
if [ ! -d "migrations" ]; then
  echo "ERROR: migrations directory not found!"
  ls -la
  exit 1
fi

# List migration files for debugging
echo "Found migration files:"
ls -la migrations/

# Run migrations with proper globbing
for file in migrations/*.sql; do
  # Check if file exists (in case glob doesn't match)
  if [ -f "$file" ]; then
    echo "Running migration: $file"
    PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f "$file"
  else
    echo "No migration files found in migrations/"
    exit 1
  fi
done

echo "Migrations complete - starting server"

# Execute the server
exec /app/server