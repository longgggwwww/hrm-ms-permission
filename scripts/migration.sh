#!/bin/bash

# Check if migration name is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <migration_name>"
  exit 1
fi

MIGRATION_NAME=$1

# Run Atlas migration diff command
docker run --rm \
  -v "$(pwd)/ent/migrate/migrations:/migrations" \
  -v "$(pwd)/ent/schema:/schema" \
  --network host \
  arigaio/atlas:latest \
  migrate diff "$MIGRATION_NAME" \
  --dir "file:///migrations" \
  --to "ent://schema" \
  --dev-url "sqlite://file?mode=memory&_fk=1"