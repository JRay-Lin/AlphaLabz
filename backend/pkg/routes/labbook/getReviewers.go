package labbook

import (
	"alphalabz/pkg/casbin"
	"alphalabz/pkg/pocketbase"
	"alphalabz/pkg/tools"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Verifier struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func GetAvailiableReviewers(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient, ce *casbin.CasbinEnforcer) {
	// Check if the request method is GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get request authorization header
	rawToken, err := tools.TokenExtractor(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Verify user permission using Casbin enforcer
	hasPermission, _, err := ce.VerifyJWTPermission(pbClient, rawToken, casbin.PermissionConfig{
		Resources: "lab_books",
		Actions:   "create",
		Scopes:    "own",
	})
	if err != nil || !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get users with permission to update lab books status
	hasPermissionList, err := ce.GetRoleIDsByPermission("lab_books", "update", "status")
	if err != nil {
		http.Error(w, "Failed to get role IDs by permission", http.StatusInternalServerError)
		return
	}

	// Format the role filter correctly using OR conditions
	roleConditions := []string{}
	for _, roleID := range hasPermissionList {
		roleConditions = append(roleConditions, fmt.Sprintf("role='%s'", roleID))
	}

	// Combine conditions with OR (||)
	rawFilter := "(" + strings.Join(roleConditions, " || ") + ")"

	// Properly encode without adding unwanted `+`
	roleFilter := strings.ReplaceAll(url.QueryEscape(rawFilter), "+", "%20")

	reviewers, _, err := pbClient.ListUsers([]string{"id", "name", "role"}, []string{}, roleFilter)
	if err != nil {
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"message": "success", "reviewers": reviewers})
}
