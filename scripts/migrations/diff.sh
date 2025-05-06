#!/bin/bash

# Add an option for migration name
migration_name=$1

# Check if migration_name is provided
if [ -z "$migration_name" ]; then
  echo "Migration name is required. Provide it as the first argument."
  exit 1
fi

# Run Atlas migration diff command
atlas migrate diff "$migration_name" \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/16-alpine/dev"