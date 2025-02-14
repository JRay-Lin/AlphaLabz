package pocketbase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// RegisterUser registers a new user in the "users" collection
func (pbClient *PocketBaseClient) NewUser(email string, password string, role string, token string) error {
	url := fmt.Sprintf("%s/api/collections/users/records", pbClient.BaseURL)

	// Data payload for the new user
	data := map[string]interface{}{
		"email":           email,
		"password":        password,
		"passwordConfirm": password,
		"role":            role,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Execute the request
	resp, err := pbClient.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to register user: non-200 status code")
	}
	return nil
}
