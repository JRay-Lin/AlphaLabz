package routes

import (
	"elimt/pkg/pocketbase"
	"encoding/json"
	"net/http"
	"strings"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Auth handlers
func HandleRegister(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient) {
	var registerData RegisterRequest
	// Get request header
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}

	// Get request body
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	role := strings.ToLower(registerData.Role)
	if role == "admin" {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	// Get available roles' id and find id that matches the role
	availableRoles, err := pbClient.GetAvailableRoles()
	var roleId string
	if err != nil {
		http.Error(w, "Failed to get available roles", http.StatusInternalServerError)
		return
	} else {
		for _, availableRole := range availableRoles {
			if strings.ToLower(availableRole.Name) == role {
				roleId = availableRole.Id
				break
			}
		}
	}

	err = pbClient.RegisterUser(registerData.Email, registerData.Password, roleId, authToken)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
