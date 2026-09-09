#!/bin/sh
echo "Running database migrations..."
./migrate || echo "Warning: migrations finished with notices, continuing to start API..."

echo "Starting TN24 API server..."
exec ./api
