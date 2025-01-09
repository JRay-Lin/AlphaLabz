#!/bin/bash
set -e

# Variables
SYSTEM="linux_amd64"
PB_BINARY="./pocketbase"
PB_VERSION="0.23.4"
PB_URL="http://localhost:8090"
PB_DOWNLOAD_URL="https://github.com/pocketbase/pocketbase/releases/download/v${PB_VERSION}/pocketbase_${PB_VERSION}_${SYSTEM}.zip"

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
    wget -q "${PB_DOWNLOAD_URL}" -O pocketbase.zip
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

# Grant superuser token
echo "Authenticating superuser..."
response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "{
        \"identity\": \"$SUPERUSER_EMAIL\",
        \"password\": \"$SUPERUSER_PASSWORD\"
    }" \
    "http://localhost:8090/api/collections/_superusers/auth-with-password")

# Extract the token from the response using jq (ensure jq is installed)
TOKEN=$(echo "$response" | grep -o '"token":"[^"]*' | sed 's/"token":"//')

# Verify if token was granted successfully
if [[ "$TOKEN" == "null" || -z "$TOKEN" ]]; then
    echo "Failed to authenticate superuser!"
    echo "Response: $response"
    exit 1
fi

# Display the granted token (for debugging purposes only, avoid this in production)
echo "Superuser token granted successfully!"
echo "Token: $TOKEN"

# Example of using the token to create a new user
echo "Creating a new user with the granted token..."

create_user_response=$(curl -s -X POST \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"$SUPERUSER_EMAIL\",
        \"password\": \"$SUPERUSER_PASSWORD\",
        \"passwordConfirm\": \"$SUPERUSER_PASSWORD\",
        \"role\": \"admin\"
    }" \
    "$PB_URL/api/collections/users/records")

# Check the result of user creation
if echo "$create_user_response" | grep -q '"code":'; then
    echo "Failed to create a new user. Response:"
    echo "$create_user_response"
    exit 1
else
    echo "New user created successfully!"
fi

# Keep PocketBase running in the foreground
wait $PB_PID