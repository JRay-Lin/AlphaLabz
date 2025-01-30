package pocketbase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// PocketBaseClient interacts with the PocketBase HTTP API
type PocketBaseClient struct {
	BaseURL string
}

// Check pocketbase connection is working
func (p *PocketBaseClient) CheckConnection() error {
	url := fmt.Sprintf("%s/api/health", p.BaseURL)

	// Create a new HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Make GET request to health endpoint
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to PocketBase: %w", err)
	}
	defer resp.Body.Close()

	// Check if status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PocketBase health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// NewPocketBaseClient initializes a new PocketBase client
func NewPocketBaseClient(baseURL string) *PocketBaseClient {
	return &PocketBaseClient{BaseURL: baseURL}
}

// RegisterUser registers a new user in the "users" collection
func (p *PocketBaseClient) RegisterUser(email string, password string, role string, authToken string) error {
	url := fmt.Sprintf("%s/api/collections/users/records", p.BaseURL)

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
	req.Header.Set("Authorization", authToken)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
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

// AuthenticateUser authenticates a user and returns their token
func (p *PocketBaseClient) AuthenticateUser(email, password string) (string, error) {
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
