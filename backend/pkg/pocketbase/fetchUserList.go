package pocketbase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// User represents a user record from PocketBase
type User struct {
	ID      string    `json:"id"`
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
func (p *PocketBaseClient) ListUsers() ([]User, error) {
	url := fmt.Sprintf("%s/api/collections/users/records", p.BaseURL)

	// Make GET request to fetch users
	resp, err := http.Get(url)
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
