package pocketbase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// User represents a user record from PocketBase
type User struct {
	Id    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	// Verified bool   `json:"verified,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	// Role    string `json:"role,omitempty"`
	Expand  expandRole `json:"expand,omitempty"`
	Gender  string     `json:"gender,omitempty"`
	Created string     `json:"created,omitempty"`
	Updated string     `json:"updated,omitempty"`
}

type expandRole struct {
	Role struct {
		Id   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}
}

// ListUsersResponse represents the PocketBase API response for listing users
type ListUsersResponse struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
	TotalItems int    `json:"totalItems"`
	Items      []User `json:"items"`
}

// ListUsers fetches the list of users
func (pbClient *PocketBaseClient) ListUsers(fields []string) (userList []User, totalUsers int, err error) {
	url := fmt.Sprintf("%s/api/collections/users/records?expand=role&", pbClient.BaseURL)

	// Add fields as query parameters if specified
	if len(fields) > 0 {
		url += "&fields=" + strings.Join(fields, ",")
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Add Authorization header
	req.Header.Add("Authorization", "Bearer "+pbClient.SuperToken)

	resp, err := pbClient.HTTPClient.Do(req)
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
