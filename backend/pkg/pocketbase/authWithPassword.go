package pocketbase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// AuthenticateUser authenticates a user and returns their token
func (p *PocketBaseClient) AuthUserWithPassword(email, password string) (string, error) {
	url := fmt.Sprintf("%s/api/collections/users/auth-with-password", p.BaseURL)

	// Data payload for authentication
	data := map[string]interface{}{
		"identity": email,
		"password": password,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	// HTTP POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to authenticate user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to authenticate user: non-200 status code")
	}

	// Parse response to extract token
	var respData struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return respData.Token, nil
}
