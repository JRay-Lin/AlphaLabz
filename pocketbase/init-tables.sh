#!/bin/bash
set -e  # Exit on error

# Token is passed as an argument
TOKEN="$1"

if [ -z "$TOKEN" ]; then
    echo "Token is required for initializing tables."
    exit 1
fi

# Role Upsert Function (Fixed for ID Conflict)
upsert_role() {
    local role_name="$1"
    local role_permission="$2"
    local role_description="$3"

    # Check if the role already exists
    EXISTING_ROLE_RESPONSE=$(curl -s -X GET "http://127.0.0.1:8090/api/collections/roles/records?filter=name='$role_name'" \
    -H "Authorization: Bearer $TOKEN")

    ROLE_ID=$(echo "$EXISTING_ROLE_RESPONSE" | grep -o '"id":"[^"]*' | cut -d '"' -f 4)

    # If the role exists, update it
    if [ -n "$ROLE_ID" ]; then
        echo "Role '$role_name' already exists. Updating..."
        curl -s -X PATCH "http://127.0.0.1:8090/api/collections/roles/records/$ROLE_ID" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "permission": '"$role_permission"',
            "description": "'"$role_description"'"
        }' | grep -q '"id":"'

        if [ $? -eq 0 ]; then
            echo "Role '$role_name' updated successfully."
        else
            echo "Failed to update role: $role_name"
            exit 1
        fi
    else
        # If the role does not exist, create it without specifying an ID
        echo "Creating new role: $role_name"
        CREATE_ROLE_RESPONSE=$(curl -s -X POST "http://127.0.0.1:8090/api/collections/roles/records" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "name": "'"$role_name"'",
            "permission": '"$role_permission"',
            "description": "'"$role_description"'"
        }')

        # Check if creation was successful
        if echo "$CREATE_ROLE_RESPONSE" | grep -q '"id":"'; then
            echo "Role '$role_name' created successfully."
        else
            echo "Failed to create role: $role_name"
            echo "$CREATE_ROLE_RESPONSE"
            exit 1
        fi
    fi
}

# Create or Update Roles
upsert_role "ADMIN" 100 "Full administrative access"
upsert_role "MODERATOR" 80 "Profesor, Teacher"
upsert_role "USER" 30 "Student"

echo "All roles have been initialized successfully."