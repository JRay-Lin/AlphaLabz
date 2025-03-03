package labbook

import (
	"alphalabz/pkg/casbin"
	"alphalabz/pkg/pocketbase"
	"alphalabz/pkg/tools"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleLabbookUploadHistory(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient, ce *casbin.CasbinEnforcer) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	rawToekn, err := tools.TokenExtractor(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check if the user has permission to view lab book upload history
	hasPermission, _, err := ce.VerifyJWTPermission(pbClient, rawToekn, casbin.PermissionConfig{
		Resources: "lab_books",
		Actions:   "view",
		Scopes:    "own"})
	if err != nil || !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	userId, err := tools.GetUserIdFromJWT(rawToekn)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Get the lab book upload history from the database
	uploadHistory, err := pbClient.ListLabbooks(fmt.Sprintf("creator='%s'", userId), []string{"*"})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get lab book upload history", http.StatusInternalServerError)
		return
	}

	// Return the lab book upload history as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(uploadHistory)
}
