#!/bin/bash
set -e

# Function to create a database if it doesn't exist
create_database() {
  local database=$1
  echo "Creating database '$database' if it doesn't exist..."
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    SELECT 'CREATE DATABASE $database'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$database');
EOSQL
}

# Create the test database
create_database 'beautix_test'

# Add extensions to both databases
for DB in beautix beautix_test; do
  echo "Setting up extensions for $DB..."
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$DB" <<-EOSQL
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";
EOSQL
done

echo "Database initialization completed!"