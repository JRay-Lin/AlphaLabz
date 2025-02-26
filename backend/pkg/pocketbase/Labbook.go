package pocketbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Labbook struct {
	Id            string   `json:"id"`
	Title         string   `json:"title,omitempty"`
	Description   string   `json:"description,omitempty"`
	Creator       string   `json:"creator,omitempty"`
	Reviewer      string   `json:"reviewer,omitempty"`
	ReviewStatus  string   `json:"review_status,omitempty"`
	ReviewComment string   `json:"review_comment,omitempty"`
	File          string   `json:"file,omitempty"`
	Attachments   []string `json:"attachments,omitempty"`
	AccessList    []string `json:"access_list,omitempty"`
	CreatedAt     string   `json:"created,omitempty"`
	UpdatedAt     string   `json:"updated,omitempty"`
}

func (pbClient *PocketBaseClient) ViewLabbook(id string, fileds []string) (Labbook, error) {
	url := fmt.Sprintf("%s/api/collections/lab_books/records/%s?fields=%s", pbClient.BaseURL, id, strings.Join(fileds, ","))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Labbook{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pbClient.SuperToken))

	resp, err := pbClient.HTTPClient.Do(req)
	if err != nil {
		return Labbook{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var labbook Labbook
	if err = json.NewDecoder(resp.Body).Decode(&labbook); err != nil {
		return Labbook{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return labbook, nil
}

func (pbClient *PocketBaseClient) UpdateLabbook(id string, data map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/collections/lab_books/records/%s", pbClient.BaseURL, id)

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pbClient.SuperToken))

	resp, err := pbClient.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
