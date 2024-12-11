#!/bin/bash
set -e

# Variables
PB_BINARY="./pocketbase"
PB_VERSION="0.23.4"
PB_URL="https://github.com/pocketbase/pocketbase/releases/download/v${PB_VERSION}/pocketbase_${PB_VERSION}_linux_amd64.zip"

# Function to wait for PocketBase to be ready
wait_for_pocketbase() {
    local retries=30
    local wait_time=2
    local endpoint="http://localhost:8090/api/health"
    
    echo "Waiting for PocketBase to be ready..."
    while [ $retries -gt 0 ]; do
        if curl -s -f "$endpoint" > /dev/null 2>&1; then
            echo "PocketBase is ready!"
            return 0
        fi
        retries=$((retries-1))
        echo "Waiting for PocketBase to start... ($retries attempts left)"
        sleep $wait_time
    done
    return 1
}

# Download PocketBase if not exists
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

# Wait for PocketBase to be ready
if ! wait_for_pocketbase; then
    echo "Failed to start PocketBase"
    kill $PB_PID
    exit 1
fi

# Create superuser account
echo "Creating or updating superuser..."
"${PB_BINARY}" superuser upsert "${ADMIN_EMAIL}" "${ADMIN_PASSWORD}"

if [ $? -eq 0 ]; then
    echo "Superuser created or updated successfully."
else
    echo "Failed to create or update superuser."
    kill $PB_PID
    exit 1
fi

# Create a temporary JSON file for the user creation request
cat > create_user.json << EOL
{
    "email": "${ADMIN_EMAIL}",
    "password": "${ADMIN_PASSWORD}",
    "passwordConfirm": "${ADMIN_PASSWORD}",
    "role": "admin"
}
EOL

# Create the user using the PocketBase API
echo "Creating new user..."
curl -X POST \
    -H "Content-Type: application/json" \
    -d @create_user.json \
    "http://localhost:8090/api/collections/users/records"

# Clean up the temporary JSON file
rm create_user.json

echo "New user has been created."

# Keep PocketBase running in the foreground
wait $PB_PID