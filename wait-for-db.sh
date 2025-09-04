#!/bin/sh

# Wait for PostgreSQL to be ready and UUID extension to be available
echo "Waiting for PostgreSQL to be ready..."

# Wait for PostgreSQL to accept connections
until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 2
done

echo "PostgreSQL is ready!"

# Wait for UUID extension to be available
echo "Waiting for UUID extension to be available..."
until psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT uuid_generate_v4();" > /dev/null 2>&1; do
  echo "UUID extension not ready - sleeping"
  sleep 2
done

echo "UUID extension is ready!"
echo "Starting SAMS backend..."
