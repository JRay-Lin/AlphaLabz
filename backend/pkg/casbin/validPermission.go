package casbin

import (
	"alphalabz/pkg/pocketbase"
	"fmt"
)

// CheckPermission validates user actions using Casbin
func (ce *CasbinEnforcer) VerifyPermission(pbClient *pocketbase.PocketBaseClient, rawJwtToken string, permissionConfig PermissionConfig) (bool, error) {
	if ce.Enforcer == nil {
		return false, fmt.Errorf("casbin Enforcer is not initialized")
	}

	userRole, err := pbClient.FetchUserInfo(rawJwtToken)
	if err != nil {
		return false, nil
	}

	ok, err := ce.Enforcer.Enforce(userRole.Role, permissionConfig.Resources, permissionConfig.Actions, permissionConfig.Scopes)
	if err != nil {
		fmt.Println("Error enforcing policy:", err)
		return false, err
	}

	return ok, nil
}
