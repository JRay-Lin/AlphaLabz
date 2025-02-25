package pocketbase

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (pbClient *PocketBaseClient) UploadLabbook(title, description, uploader, reviewer, labbookPath string, attachmentPaths []string) error {
	url := fmt.Sprintf("%s/api/collections/lab_books/records", pbClient.BaseURL)

	// Create multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Open and add labbook file
	labbook, err := os.Open(labbookPath)
	if err != nil {
		return fmt.Errorf("failed to open labbook file: %w", err)
	}
	defer labbook.Close()

	labbookPart, err := writer.CreateFormFile("file", filepath.Base(labbookPath))
	if err != nil {
		return fmt.Errorf("failed to create form file for labbook: %w", err)
	}

	_, err = io.Copy(labbookPart, labbook)
	if err != nil {
		return fmt.Errorf("failed to copy labbook content: %w", err)
	}

	// Add multiple attachment files
	for _, attachmentPath := range attachmentPaths {
		attachment, err := os.Open(attachmentPath)
		if err != nil {
			return fmt.Errorf("failed to open attachment %s: %w", attachmentPath, err)
		}
		defer attachment.Close()

		// Use the same field name "attachments" for all files
		attachmentPart, err := writer.CreateFormFile("attachments", filepath.Base(attachmentPath))
		if err != nil {
			return fmt.Errorf("failed to create form file for attachment: %w", err)
		}

		_, err = io.Copy(attachmentPart, attachment)
		if err != nil {
			return fmt.Errorf("failed to copy attachment content: %w", err)
		}
	}

	// Add form fields **AFTER** adding the files
	_ = writer.WriteField("title", title)
	_ = writer.WriteField("creator", uploader)
	_ = writer.WriteField("reviewer", reviewer)
	_ = writer.WriteField("review_status", "pending")
	_ = writer.WriteField("access_list", fmt.Sprintf("[\"%s\", \"%s\"]", uploader, reviewer))

	if description != "" {
		_ = writer.WriteField("description", description)
	}

	// Close writer to finalize form data
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Correctly set content type with the generated boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pbClient.SuperToken))

	// Send request
	resp, err := pbClient.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
