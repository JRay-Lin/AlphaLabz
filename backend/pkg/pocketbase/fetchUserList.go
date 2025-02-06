package pocketbase

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// User represents a user record from PocketBase
type User struct {
	Id              string `json:"id"`
	Email           string `json:"email"`
	EmailVisibility bool   `json:"emailVisibility"`
	Verified        bool   `json:"verified"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`
	Role            string `json:"role"`
	Gender          string `json:"gender"`
	Created         string `json:"created"`
	Updated         string `json:"updated"`
}

// ListUsersResponse represents the PocketBase API response for listing users
type ListUsersResponse struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
	TotalItems int    `json:"totalItems"`
	Items      []User `json:"items"`
}

// ListUsers fetches the list of users
func (p *PocketBaseClient) ListUsers(token string) (userList []User, totalUsers int, err error) {
	url := fmt.Sprintf("%s/api/collections/users/records", p.BaseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Add Authorization header
	req.Header.Add("Authorization", "Bearer "+token)

	// Make the request using http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("failed to fetch users: server returned status %d", resp.StatusCode)
	}

	// Parse response
	var UserListResponse ListUsersResponse
	err = json.NewDecoder(resp.Body).Decode(&UserListResponse)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return UserListResponse.Items, UserListResponse.TotalItems, nil
}
