package routes

import (
	"elimt/pkg/pocketbase"
	"encoding/json"
	"net/http"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetTimestamp() {

}

func HandleLogin(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient) {
	var loginData LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := pbClient.AuthenticateUser(loginData.Email, loginData.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	timestamp := time.Now().Unix()
	formattedTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"token":     token,
		"timestamp": formattedTime,
	})
}
