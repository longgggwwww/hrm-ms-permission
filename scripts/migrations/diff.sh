#!/bin/bash

# Add an option for migration name
while getopts n: flag
  do
    case "${flag}" in
        n) migration_name=${OPTARG};;
    esac
done

# Check if migration_name is provided
if [ -z "$migration_name" ]; then
  echo "Migration name is required. Use -n to specify the migration name."
  exit 1
fi

# Run Atlas migration diff command
atlas migrate diff "$migration_name" \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"