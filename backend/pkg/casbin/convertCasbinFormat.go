package casbin

import (
	"fmt"
	"strings"
)

// ConvertCasbinFormat converts RolePermissions into Casbin policy rules
func ConvertCasbinFormat(permissions []RolePermission) []string {
	var casbinPolicies []string

	for _, role := range permissions {
		roleId := role.Id // Role ID as Casbin "sub"

		for resource, actions := range role.Permission {
			actionList, ok := actions.([]interface{})
			if !ok {
				continue
			}

			for _, action := range actionList {
				actionStr, ok := action.(string)
				if !ok {
					continue
				}

				actionParts := strings.Split(actionStr, ":") // Split into action & scope
				actionType := actionParts[0]                 // "view", "update", etc.
				scope := "all"                               // Default scope

				if len(actionParts) > 1 {
					scopes := strings.Split(actionParts[1], ",") // Handle multiple scopes

					for _, s := range scopes {
						policy := fmt.Sprintf("p, %s, %s, %s, %s", roleId, resource, actionType, strings.TrimSpace(s))
						casbinPolicies = append(casbinPolicies, policy)
					}
					continue
				}

				// If there's only one scope, add normally
				policy := fmt.Sprintf("p, %s, %s, %s, %s", roleId, resource, actionType, scope)
				casbinPolicies = append(casbinPolicies, policy)
			}
		}
	}

	return casbinPolicies
}
