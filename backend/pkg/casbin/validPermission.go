package casbin

import (
	"alphalabz/pkg/pocketbase"
	"fmt"
)

// CheckPermission validates user actions using Casbin
//
// Require Resources, Actions and Scopes to be provided in PermissionConfig
func (ce *CasbinEnforcer) VerifyJWTPermission(pbClient *pocketbase.PocketBaseClient, rawJwtToken string, permissionConfig PermissionConfig) (reqPermission bool, starPermission bool, err error) {
	if ce.Enforcer == nil {
		return false, false, fmt.Errorf("casbin Enforcer is not initialized")
	}

	userRole, err := pbClient.ViewUser(rawJwtToken)
	if err != nil {
		return false, false, nil
	}

	// Check if the user has the '*' scope (unrestricted access).
	starScopeCheck, err := ce.Enforcer.Enforce(userRole.RoleId, permissionConfig.Resources, permissionConfig.Actions, "*")
	if err != nil {
		fmt.Println("Error enforcing policy:", err)
		return false, false, err
	}

	// Check permission using the specified scope.
	reqPermissionCheck, err := ce.Enforcer.Enforce(userRole.RoleId, permissionConfig.Resources, permissionConfig.Actions, permissionConfig.Scopes)
	if err != nil {
		fmt.Println("Error enforcing policy:", err)
		return false, false, err
	}

	return starScopeCheck, reqPermissionCheck, nil
}
