package casbin

import (
	"fmt"
)

// CheckPermission validates user actions using Casbin
func (ce *CasbinEnforcer) VerifyPermission(roleId string, permissionConfig PermissionConfig) (bool, error) {
	if ce.Enforcer == nil {
		return false, fmt.Errorf("casbin Enforcer is not initialized")
	}

	ok, err := ce.Enforcer.Enforce(roleId, permissionConfig.Resources, permissionConfig.Actions, permissionConfig.Scopes)
	if err != nil {
		fmt.Println("Error enforcing policy:", err)
		return false, err
	}
	return ok, nil
}
