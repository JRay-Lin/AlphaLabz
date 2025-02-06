package casbin

import (
	"alphalabz/pkg/pocketbase"
	"fmt"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// CasbinEnforcer struct to manage RBAC enforcement
type CasbinEnforcer struct {
	Enforcer *casbin.Enforcer
}

type PermissionRespond struct {
	Items      []RolePermission `json:"items"`
	Page       int              `json:"page"`
	PerPage    int              `json:"perPage"`
	TotalItems int              `json:"totalItems"`
	TotalPages int              `json:"totalPages"`
}

type RolePermission struct {
	// CollectionId   string                 `json:"collectionId"`
	// CollectionName string                 `json:"collectionName"`
	// Description string                 `json:"description"`
	Id string `json:"id"`
	// Name        string                 `json:"name"`
	Permission map[string]interface{} `json:"permission"`
}

// InitializeCasbin initializes Casbin and loads permissions dynamically
func InitializeCasbin(pbClient *pocketbase.PocketBaseClient) (*CasbinEnforcer, error) {
	// Ensure policy file exists
	policyFile := "casbin_policies.csv"
	if _, err := os.Stat(policyFile); os.IsNotExist(err) {
		// Create an empty policy file
		file, err := os.Create(policyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create Casbin policy file: %v", err)
		}
		file.Close()
	}

	// Define RBAC model
	rbacModel := `
		[request_definition]
		r = sub, obj, act, scope

		[policy_definition]
		p = sub, obj, act, scope

		[role_definition]
		g = _, _

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act && r.scope == p.scope
	`

	// Create model and adapter
	m, err := model.NewModelFromString(rbacModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin model: %v", err)
	}

	adapter := fileadapter.NewAdapter(policyFile)

	// Create a new Casbin enforcer
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create Casbin enforcer: %v", err)
	}

	// Load existing policies (if any)
	err = e.LoadPolicy()
	if err != nil {
		fmt.Println("No existing policies found, continuing with an empty policy set.")
	}

	// Fetch permissions from PocketBase
	permissions := FetchPermissions(pbClient)
	if permissions == nil {
		return nil, fmt.Errorf("failed to fetch permissions from PocketBase")
	}

	// Convert permissions to Casbin format and save
	policies := ConvertCasbinFormat(permissions)
	SavePoliciesToFile(policies, policyFile)

	// Reload Casbin policies after update
	err = e.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to reload Casbin policies: %v", err)
	}

	fmt.Println("Casbin Enforcer initialized successfully.")
	return &CasbinEnforcer{Enforcer: e}, nil
}

// CheckPermission validates user actions using Casbin
func (ce *CasbinEnforcer) CheckPermission(roleId, resource, action, scope string) bool {
	ok, err := ce.Enforcer.Enforce(roleId, resource, action, scope)
	if err != nil {
		fmt.Println("Error enforcing policy:", err)
		return false
	}
	return ok
}
