package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CleanUploads cleans up old or unnecessary uploads.
func CleanUploads() error {
	folder := "./uploads"

	// Delete all files in the folder
	if err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fmt.Println("Deleting file:", path)
			os.Remove(path)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func StartAutoCleanUploads(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			err := CleanUploads()
			if err != nil {
				fmt.Println("Error cleaing uploads:", err)
			}
		}
	}()
}

// CreateUploadsDir creates the necessary directories for file uploads.
func CreateUploadsDir() error {
	// Define the directories to be created
	directories := []string{
		"./uploads/avatar",
		"./uploads/labbook",
	}

	// Loop through and create each directory
	for _, dir := range directories {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	return nil
}
