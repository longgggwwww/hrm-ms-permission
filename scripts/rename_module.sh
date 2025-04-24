#!/bin/bash

# This script renames the Go module in the go.mod file and updates all import paths in the project.

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <new-module-name>"
  exit 1
fi

NEW_MODULE_NAME=$1

# Update the module name in go.mod
sed -i "s/^module .*/module $NEW_MODULE_NAME/" go.mod

# Find and replace all import paths in the project
grep -rl "$(grep '^module ' go.mod | cut -d' ' -f2)" . | grep -v "\.git" | xargs sed -i "s|$(grep '^module ' go.mod | cut -d' ' -f2)|$NEW_MODULE_NAME|g"

echo "Module name updated to $NEW_MODULE_NAME and import paths updated."