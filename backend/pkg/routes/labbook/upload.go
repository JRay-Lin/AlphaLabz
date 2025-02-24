package labbook

import (
	"alphalabz/pkg/casbin"
	"alphalabz/pkg/pocketbase"
	"net/http"
)

func HandleLabBookUpload(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient, ce *casbin.CasbinEnforcer) {
	var permissionConfig = casbin.PermissionConfig{
		Resources: "labbook",
		Actions:   "create",
		Scopes:    "own",
	}

	// Check if the user has permission to update lab books
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	scopes, err := ce.ScopeFetcher(pbClient, r.Header.Get("Authorization"), permissionConfig)
	if err != nil {
		http.Error(w, "Failed to fetch scopes", http.StatusInternalServerError)
		return
	}

	if len(scopes) == 0 {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

}
