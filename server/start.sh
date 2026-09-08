#!/bin/sh
set -e

echo "Running database migrations..."
./migrate

echo "Starting TN24 API server..."
exec ./api
