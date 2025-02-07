package user

import (
	"net/http"
)

type NewUserRequest struct {
	Token     string  `json:"token"`
	Timestamp string  `json:"timestamp"`
	NewUser   NewUser `json:"newUser"`
}

type NewUser struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func HandleInviteNewUser(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}
