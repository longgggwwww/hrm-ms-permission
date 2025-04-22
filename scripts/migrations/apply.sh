#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Check if DB_URL environment variable is set
if [ -z "$DB_URL" ]; then
  echo "Error: DB_URL environment variable is not set."
  exit 1
fi

# Apply migrations using Atlas
atlas migrate apply \
  --dir "file://ent/migrate/migrations" \
  --url "$DB_URL"
  