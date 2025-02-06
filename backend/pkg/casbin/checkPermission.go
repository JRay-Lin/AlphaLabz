package casbin

import (
	"fmt"
)

// CheckPermission validates user actions using Casbin
func (ce *CasbinEnforcer) CheckPermission(roleId, resource, action, scope string) bool {
	ok, err := ce.Enforcer.Enforce(roleId, resource, action, scope)
	if err != nil {
		fmt.Println("Error enforcing policy:", err)
		return false
	}
	return ok
}
