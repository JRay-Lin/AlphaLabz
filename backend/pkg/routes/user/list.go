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

// List all users in the database
//
// ✅ Authorization:
// - Requires an `Authorization` header with a valid token.
//
// ✅ Successful Response (200 OK):
//
//		{
//	    "totalUsers": 1,
//	    "users": [
//	        {
//	            "id": "qz73n36tig1k7z7",
//	            "email": "test@alphalabz.net",
//	            "emailVisibility": false,
//	            "verified": false,
//	            "name": "",
//	            "avatar": "test.png",
//	            "role": "ADMIN",
//	            "gender": "",
//	            "created": "2025-01-14 12:35:58.273Z",
//	            "updated": "2025-01-30 10:00:14.643Z"
//	        },
//			...
//	    ]
//	}
//
// ❌ Error Responses:
//   - 401 Unauthorized → Missing or Invalid Authorization token
//   - 500 Internal Server Error → Server issue
func HandleUserList(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get token from request header
	token := r.Header.Get("Authorization")

	// Remove "Bearer " prefix from token
	rawToken := token[len("Bearer"):]
	if rawToken == "" {
		http.Error(w, "Missing Authorization token", http.StatusUnauthorized)
		return
	}

	// Fetch users from PocketBase
	users, totalUsers, err := pbClient.ListUsers(token)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Invalid token
	if totalUsers == 0 {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
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

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
