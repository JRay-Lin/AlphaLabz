package labbook

import (
	"alphalabz/pkg/casbin"
	"alphalabz/pkg/pocketbase"
	"net/http"
)

func HandleLabBookUpdate(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient, ce *casbin.CasbinEnforcer) {
	// var permissionConfig = casbin.PermissionConfig{
	// 	Resources: "labbook",
	// 	Actions:   "update",
	// }

	// // Check if the user has permission to update lab books
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	// scopes, err := tools.ScopeFetcher(pbClient, ce, r.Header.Get("Authorization"), permissionConfig.Resources, permissionConfig.Actions)
	// if err != nil {
	// 	http.Error(w, "Failed to fetch scopes", http.StatusInternalServerError)
	// 	return
	// }

}
