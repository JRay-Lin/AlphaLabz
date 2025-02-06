package pocketbase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Define a struct to parse the verification response
type VerifiedResult struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"userId,omitempty"`
	Role   string `json:"role,omitempty"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
}

type pbRespond struct {
	Record pbUser `json:"record"`
	Token  string `json:"token"`
}

type pbUser struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// VerifyToken sends a request to check if the token is valid.
func (p *PocketBaseClient) VerifyToken(token string) (VerifiedResult, error) {
	url := fmt.Sprintf("%s/api/collections/users/auth-refresh", p.BaseURL)

	// Create a new request with the POST method
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return VerifiedResult{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Add Authorization header
	req.Header.Add("Authorization", "Bearer "+token)

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return VerifiedResult{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read raw response body for debugging
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return VerifiedResult{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-OK status codes
	if resp.StatusCode != http.StatusOK {
		return VerifiedResult{Valid: false}, fmt.Errorf("token verification failed: %s", string(rawBody))
	}

	// Parse response JSON into struct
	var pbRes pbRespond
	if err := json.Unmarshal(rawBody, &pbRes); err != nil {
		return VerifiedResult{}, fmt.Errorf("failed to parse response: %w", err)
	}

	// Return parsed result
	return VerifiedResult{
		Valid:  true,
		UserID: pbRes.Record.Id,
		Role:   pbRes.Record.Role,
		Name:   pbRes.Record.Name,
	}, nil
}
