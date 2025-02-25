package pocketbase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// User represents a user record from PocketBase
type User struct {
	Id        string  `json:"id,omitempty"`
	Email     string  `json:"email,omitempty"`
	Name      string  `json:"name,omitempty"`
	Avatar    string  `json:"avatar,omitempty"`
	Gender    string  `json:"gender,omitempty"`
	RoleId    string  `json:"role,omitempty"`
	SettingId string  `json:"settings,omitempty"`
	BirthDate string  `json:"birthdate,omitempty"`
	Expand    *Expand `json:"expand,omitempty"`
	Created   string  `json:"created,omitempty"`
	Updated   string  `json:"updated,omitempty"`
}

type Expand struct {
	Role        *Role        `json:"role,omitempty"`
	UserSetting *UserSetting `json:"userSetting,omitempty"`
}
type UserSetting struct {
	Id          string `json:"id,omitempty"`
	AppLanguage string `json:"app_language,omitempty"`
	Theme       string `json:"theme,omitempty"`
}

// ListUsersResponse represents the PocketBase API response for listing users
type ListUsersResponse struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
	TotalItems int    `json:"totalItems"`
	Items      []User `json:"items"`
}

// ListUsers fetches the list of users
func (pbClient *PocketBaseClient) ListUsers(fields []string, expand []string, filter string) (userList []User, totalUsers int, err error) {
	url := fmt.Sprintf("%s/api/collections/users/records", pbClient.BaseURL)

	// Add fields as query parameters if specified
	url += "?fields=" + strings.Join(fields, ",")

	if len(expand) != 0 {
		url += "&expand=" + strings.Join(expand, ",")
	}

	if filter != "" {
		url += "&filter=" + filter
	}

	// Create HTTP request
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
	var userListResp ListUsersResponse
	err = json.NewDecoder(resp.Body).Decode(&userListResp)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return userListResp.Items, userListResp.TotalItems, nil
}
