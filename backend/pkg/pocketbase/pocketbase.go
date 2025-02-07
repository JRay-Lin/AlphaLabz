package pocketbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// PocketBaseClient interacts with the PocketBase HTTP API
type PocketBaseClient struct {
	BaseURL    string
	SuperToken string
	HTTPClient *http.Client
}

// NewPocketBase initializes a new PocketBase client, authenticates, and verifies the connection.
func NewPocketBase(baseURL, superuserEmail, superuserPassword string, maxRetries int, retryInterval time.Duration) (*PocketBaseClient, error) {
	client := &PocketBaseClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second}}

	// Authenticate superuser and store the token
	token, err := client.authenticateSuperuser(superuserEmail, superuserPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate superuser: %w", err)
	}
	log.Println("Successfully authenticated superuser")
	client.SuperToken = token

	// Verify PocketBase connection with retries
	for i := 0; i < maxRetries; i++ {
		err := client.CheckConnection()
		if err == nil {
			log.Println("Successfully connected to PocketBase")
			return client, nil
		}
		log.Printf("Failed to connect to PocketBase, attempt %d/%d. Retrying in %s...", i+1, maxRetries, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to connect to PocketBase after %d attempts", maxRetries)
}

// authenticateSuperuser logs in the superuser and retrieves the authentication token
func (pbClient *PocketBaseClient) authenticateSuperuser(email, password string) (string, error) {
	url := fmt.Sprintf("%s/api/collections/_superusers/auth-with-password", pbClient.BaseURL)

	// Data payload for authentication
	data := map[string]interface{}{
		"identity": email,
		"password": password,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := pbClient.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate superuser: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to authenticate superuser: status %d", resp.StatusCode)
	}

	// Parse response
	var respData struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return respData.Token, nil
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
