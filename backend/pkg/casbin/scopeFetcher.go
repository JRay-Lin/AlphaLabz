package casbin

import (
	"alphalabz/pkg/pocketbase"
)

func (ce *CasbinEnforcer) ScopeFetcher(pbClient *pocketbase.PocketBaseClient, rawJwtToken string, permissionCfg PermissionConfig) (scopes []string, err error) {
	userRole, err := pbClient.FetchUserInfo(rawJwtToken)
	if err != nil {
		return scopes, nil
	}

	scopes, err = ce.CheckPermissionScopes(userRole.Role, permissionCfg.Resources, permissionCfg.Actions)
	if err != nil {
		return scopes, err
	}

	return scopes, nil
}
