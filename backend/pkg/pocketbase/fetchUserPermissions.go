package pocketbase

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PermissionData map[string]map[string][]string

func (p *PocketBaseClient) FetchUserPermissions(userJwt string) (PermissionData, error) {
	url := fmt.Sprintf("%s/api/permissions/%s", p.BaseURL, userJwt)

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add autorization header
	req.Header.Set("Authorization", "Bearer "+p.SuperToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var userPermissionResponse PermissionData
	if err := json.NewDecoder(resp.Body).Decode(&userPermissionResponse); err != nil {
		return nil, err
	}

	return userPermissionResponse, nil
}
