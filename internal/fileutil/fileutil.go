package fileutil

import (
	"fmt"
	"os"
)

// WriteChangelogToFile writes the changelog content to a specified file.
func WriteChangelogToFile(filePath, content string) error {
	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists: %s", filePath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking file existence: %w", err)
	}

	// Create the file and write content
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
