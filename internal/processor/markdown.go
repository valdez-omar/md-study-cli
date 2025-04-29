package processor

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/valdezdata/md-study/internal/storage"
)

// ImportMarkdownFiles processes all markdown files in a directory
func ImportMarkdownFiles(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".md") {
			filePath := filepath.Join(dirPath, file.Name())
			if err := processMarkdownFile(filePath); err != nil {
				return fmt.Errorf("failed to process file %s: %w", file.Name(), err)
			}
		}
	}

	return nil
}

// processMarkdownFile reads a markdown file and extracts content for flashcards
func processMarkdownFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// For now, just store the raw content
	// In a more advanced version, you'd parse the markdown and extract key concepts
	note := storage.Note{
		FilePath:   filePath,
		Filename:   filepath.Base(filePath),
		RawContent: string(content),
		LastImport: time.Now(),
	}

	return storage.SaveNote(note)
}
