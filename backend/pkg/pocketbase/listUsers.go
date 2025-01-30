package pocketbase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// User represents a user record from PocketBase
type User struct {
	Id      string    `json:"id"`
	Email   string    `json:"email"`
	Role    string    `json:"role"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// ListUsersResponse represents the PocketBase API response for listing users
type ListUsersResponse struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
	TotalItems int    `json:"totalItems"`
	Items      []User `json:"items"`
}

// ListUsers fetches all users from the PocketBase users collection
func (p *PocketBaseClient) ListUsers(token string) ([]User, error) {
	url := fmt.Sprintf("%s/api/collections/users/records", p.BaseURL)

	// Create a new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add Authorization header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// Make the request using http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch users: server returned status %d", resp.StatusCode)
	}

	// Parse response
	var listResp ListUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Items, nil
}
