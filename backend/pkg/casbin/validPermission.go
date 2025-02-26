package casbin

import (
	"alphalabz/pkg/pocketbase"
	"fmt"
)

// CheckPermission validates user actions using Casbin
//
// Require Resources, Actions and Scopes to be provided in PermissionConfig
func (ce *CasbinEnforcer) VerifyJWTPermission(pbClient *pocketbase.PocketBaseClient, rawJwtToken string, permissionConfig PermissionConfig) (bool, error) {
	if ce.Enforcer == nil {
		return false, fmt.Errorf("casbin Enforcer is not initialized")
	}

	userRole, err := pbClient.ViewUser(rawJwtToken)
	if err != nil {
		return false, nil
	}

	// Check if the user has the '*' scope (unrestricted access).
	starScopeCheck, err := ce.Enforcer.Enforce(userRole.RoleId, permissionConfig.Resources, permissionConfig.Actions, "*")
	if err != nil {
		fmt.Println("Error enforcing policy:", err)
		return false, err
	}
	if starScopeCheck {
		return true, nil
	}

	// Check permission using the specified scope.
	ok, err := ce.Enforcer.Enforce(userRole.RoleId, permissionConfig.Resources, permissionConfig.Actions, permissionConfig.Scopes)
	if err != nil {
		fmt.Println("Error enforcing policy:", err)
		return false, err
	}

	return ok, nil
}
