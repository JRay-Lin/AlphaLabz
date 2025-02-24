package user

import (
	"alphalabz/pkg/casbin"
	"alphalabz/pkg/pocketbase"
	"encoding/json"
	"net/http"
)

func HandleUserRemove(w http.ResponseWriter, r *http.Request, pbClient *pocketbase.PocketBaseClient, ce *casbin.CasbinEnforcer) {
	// var permissionConfig = casbin.PermissionConfig{
	// 	Resources: "user",
	// 	Actions:   "delete",
	// }

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Response
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
