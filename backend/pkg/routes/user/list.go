package user

import (
	"alphalabz/pkg/pocketbase"
	"encoding/json"
	"net/http"
)

// UserListResponse represents the response structure for the user list endpoint
type UserListResponse struct {
	TotalUsers int               `json:"totalUsers"`
	Users      []pocketbase.User `json:"users"`
}

// HandleUserList handles the GET request to fetch all users
func HandleUserList(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient) {
	// Get token from request header
	token := r.Header.Get("Authorization")

	// Fetch users from PocketBase
	users, totalUsers, err := pbClient.ListUsers(token)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Fetch available roles
	availableRoles, err := pbClient.GetAvailableRoles()
	if err != nil {
		http.Error(w, "Failed to get available roles", http.StatusInternalServerError)
		return
	}

	// Create a map to store role ID -> role name
	roleMap := make(map[string]string)
	for _, role := range availableRoles {
		roleMap[role.Id] = role.Name
	}

	// Replace role IDs in users with role names
	for i, user := range users {
		if roleName, exists := roleMap[user.Role]; exists {
			users[i].Role = roleName // Replace role ID with role name
		} else {
			users[i].Role = "UNKNOWN" // Fallback if role ID is not found
		}
	}

	// Prepare response
	response := UserListResponse{
		TotalUsers: totalUsers,
		Users:      users,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Encode and send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
