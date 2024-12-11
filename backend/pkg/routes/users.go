package routes

import (
	"elimt/pkg/pocketbase"
	"encoding/json"
	"net/http"
)

// UserListResponse represents the response structure for the user list endpoint
type UserListResponse struct {
	Users []pocketbase.User `json:"users"`
}

// HandleUserList handles the GET request to fetch all users
func HandleUserList(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get token from Authorization header
	token, err := ExtractAuthToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch users from PocketBase
	users, err := pbClient.ListUsers(token)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := UserListResponse{
		Users: users,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Encode and send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
