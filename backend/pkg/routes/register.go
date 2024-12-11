package routes

import (
	"elimt/pkg/pocketbase"
	"encoding/json"
	"net/http"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Auth handlers
func HandleRegister(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient) {
	var registerData RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if registerData.Role != "moderator" && registerData.Role != "user" {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	err := pbClient.RegisterUser(registerData.Email, registerData.Password, registerData.Role)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
