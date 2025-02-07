package user

import (
	"alphalabz/pkg/casbin"
	"alphalabz/pkg/pocketbase"
	"alphalabz/pkg/tools"
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
func HandleUserList(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient, ce *casbin.CasbinEnforcer) {
	var permissionConfig = casbin.PermissionConfig{
		Resources: "users",
		Actions:   "list",
	}

	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get token from request header
	rawJwtToken, err := tools.TokenExtractor(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userRole, err := pbClient.FetchUserRole(rawJwtToken)
	if err != nil {
		http.Error(w, "Failed to fetch user role", http.StatusInternalServerError)
		return
	}

	var reqFields []string
	scopes, err := ce.CheckPermissionScopes(userRole.RoleId, permissionConfig.Resources, permissionConfig.Actions)
	if err != nil {
		http.Error(w, "Failed to check permission", http.StatusInternalServerError)
	} else {
		// Check if user has "all" scope
		for _, scope := range scopes {
			if scope == "all" {
				reqFields = []string{"*"} // Grant access to all fields
				break
			}
		}

		// If "all" is NOT found, use the allowed scopes
		if len(reqFields) == 0 {
			reqFields = scopes
		}
	}

	userList, TotalUsers, err := pbClient.ListUsers(reqFields)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	result := UserListResponse{
		TotalUsers: TotalUsers,
		Users:      userList,
	}

	// Encode and send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
