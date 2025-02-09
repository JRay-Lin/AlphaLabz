package pocketbase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Role struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions int8   `json:"permission"` // Fixed from "permissions" to match JSON key
}

type rolesResponse struct {
	Items []Role `json:"items"` // Extract the "items" array from response
}

// Get all available roles in the database
func (pbClient *PocketBaseClient) GetAvailableRoles(fields []string) (rolesResponse, error) {
	url := fmt.Sprintf("%s/api/collections/roles/records", pbClient.BaseURL)
	// Add fields as query parameters if specified
	url += "?fields=" + strings.Join(fields, ",")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return rolesResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := pbClient.HTTPClient.Do(req)
	if err != nil {
		return rolesResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close() // Ensure body is closed after reading

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return rolesResponse{}, fmt.Errorf("failed to fetch roles: received status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rolesResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println(string(body))

	// Parse JSON response
	var roleListResp rolesResponse
	err = json.NewDecoder(resp.Body.item).Encode(&roleListResp)
	if err != nil {
		return rolesResponse{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return roleListResp, nil
}
