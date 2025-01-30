package pocketbase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Role struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions int8   `json:"permission"` // Fixed from "permissions" to match JSON key
}

type roleResponse struct {
	Items []Role `json:"items"` // Extract the "items" array from response
}

// Get all available roles in the database
func (p *PocketBaseClient) GetAvailableRoles() ([]Role, error) {
	url := fmt.Sprintf("%s/api/collections/roles/records", p.BaseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close() // Ensure body is closed after reading

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch roles: received status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var result roleResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return result.Items, nil
}
