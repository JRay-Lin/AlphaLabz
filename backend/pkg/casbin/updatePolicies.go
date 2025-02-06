package casbin

import (
	"alphalabz/pkg/pocketbase"
	"bufio"
	"fmt"
	"os"
)

// UpdatePolicies fetches the latest permissions and replaces the old policy file
func (ce *CasbinEnforcer) UpdatePolicies(pbClient *pocketbase.PocketBaseClient) error {
	fmt.Println("Fetching latest permissions from PocketBase...")

	// Fetch latest permissions
	permissions := FetchPermissions(pbClient)
	if permissions == nil {
		return fmt.Errorf("failed to fetch latest permissions")
	}

	// Convert permissions to Casbin format
	policies := ConvertCasbinFormat(permissions)

	// Overwrite the old policy file
	policyFile := "casbin_policies.csv"
	err := overwritePolicyFile(policyFile, policies)
	if err != nil {
		return fmt.Errorf("failed to update policy file: %v", err)
	}

	// Reload policies into Casbin
	err = ce.Enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("failed to reload Casbin policies: %v", err)
	}

	fmt.Println("Casbin policies successfully updated and reloaded.")
	return nil
}

// overwritePolicyFile replaces the existing policy file with new policies
func overwritePolicyFile(filename string, policies []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create policy file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, policy := range policies {
		_, err := writer.WriteString(policy + "\n")
		if err != nil {
			return fmt.Errorf("failed to write policy: %v", err)
		}
	}
	writer.Flush()
	fmt.Println("Policy file updated successfully.")
	return nil
}
