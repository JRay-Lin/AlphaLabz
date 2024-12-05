#!/bin/bash

set -e

# Variables
PB_BINARY="./pocketbase"
PB_VERSION="0.23.4"
PB_URL="https://github.com/pocketbase/pocketbase/releases/download/v${PB_VERSION}/pocketbase_${PB_VERSION}_linux_amd64.zip"

echo "Starting PocketBase initialization..."

# Check if PocketBase binary exists
if [ ! -f "${PB_BINARY}" ]; then
    echo "PocketBase binary not found. Downloading..."
    wget -q "${PB_URL}" -O pocketbase.zip
    unzip pocketbase.zip
    rm pocketbase.zip
    chmod +x "${PB_BINARY}"
    echo "PocketBase downloaded and prepared."
else
    echo "PocketBase binary already exists. Skipping download."
fi

# Start PocketBase in the background
"${PB_BINARY}" serve --http=0.0.0.0:8090 &
PB_PID=$!

# Wait for PocketBase to start
echo "Waiting for PocketBase to start..."
sleep 5

# Create superuser account
echo "Creating or updating superuser..."
"${PB_BINARY}" superuser upsert "${ADMIN_EMAIL}" "${ADMIN_PASSWORD}"


# Check if the command was successful
if [ $? -eq 0 ]; then
    echo "Superuser created or updated successfully."
else
    echo "Failed to create or update superuser."
    kill $PB_PID
    exit 1
fi

# Wait for PocketBase to terminate
echo "PocketBase is running. Waiting for process to terminate..."
wait $PB_PID