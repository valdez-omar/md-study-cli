package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const (
	dataDir   = ".md-study"
	notesFile = "notes.json"
	cardsFile = "flashcards.json"
	statsFile = "stats.json"
)

// Initialize creates the storage directory if it doesn't exist
func Initialize() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create data directory
	storageDir := filepath.Join(homeDir, dataDir)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Initialize files if they don't exist
	for _, file := range []string{notesFile, cardsFile, statsFile} {
		path := filepath.Join(storageDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.WriteFile(path, []byte("[]"), 0644); err != nil {
				return fmt.Errorf("failed to initialize %s: %w", file, err)
			}
		}
	}

	return nil
}

// getFilePath returns the full path to a storage file
func getFilePath(filename string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, dataDir, filename), nil
}

// SaveNote saves a note to storage
func SaveNote(note Note) error {
	if err := Initialize(); err != nil {
		return err
	}

	// Generate ID if not already set
	if note.ID == "" {
		note.ID = uuid.New().String()
	}

	filePath, err := getFilePath(notesFile)
	if err != nil {
		return err
	}

	// Read existing notes
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read notes file: %w", err)
	}

	var notes []Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return fmt.Errorf("failed to parse notes: %w", err)
	}

	// Check if note exists and update or add
	found := false
	for i, n := range notes {
		if n.FilePath == note.FilePath {
			notes[i] = note
			found = true
			break
		}
	}

	if !found {
		notes = append(notes, note)
	}

	// Save back to file
	updatedData, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal notes: %w", err)
	}

	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write notes file: %w", err)
	}

	return nil
}

// GetNote retrieves a note by ID
func GetNote(id string) (Note, error) {
	filePath, err := getFilePath(notesFile)
	if err != nil {
		return Note{}, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return Note{}, fmt.Errorf("failed to read notes file: %w", err)
	}

	var notes []Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return Note{}, fmt.Errorf("failed to parse notes: %w", err)
	}

	for _, note := range notes {
		if note.ID == id {
			return note, nil
		}
	}

	return Note{}, fmt.Errorf("note not found: %s", id)
}

// GetAllNotes retrieves all notes
func GetAllNotes() ([]Note, error) {
	filePath, err := getFilePath(notesFile)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read notes file: %w", err)
	}

	var notes []Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return nil, fmt.Errorf("failed to parse notes: %w", err)
	}

	return notes, nil
}

// SaveFlashcard saves a flashcard to storage
func SaveFlashcard(card Flashcard) error {
	if err := Initialize(); err != nil {
		return err
	}

	// Generate ID if not already set
	if card.ID == "" {
		card.ID = uuid.New().String()
	}

	filePath, err := getFilePath(cardsFile)
	if err != nil {
		return err
	}

	// Read existing flashcards
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read flashcards file: %w", err)
	}

	var cards []Flashcard
	if err := json.Unmarshal(data, &cards); err != nil {
		return fmt.Errorf("failed to parse flashcards: %w", err)
	}

	// Check if card exists and update or add
	found := false
	for i, c := range cards {
		if c.ID == card.ID {
			cards[i] = card
			found = true
			break
		}
	}

	if !found {
		cards = append(cards, card)
	}

	// Save back to file
	updatedData, err := json.MarshalIndent(cards, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal flashcards: %w", err)
	}

	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write flashcards file: %w", err)
	}

	return nil
}

// GetFlashcard retrieves a flashcard by ID
func GetFlashcard(id string) (Flashcard, error) {
	filePath, err := getFilePath(cardsFile)
	if err != nil {
		return Flashcard{}, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return Flashcard{}, fmt.Errorf("failed to read flashcards file: %w", err)
	}

	var cards []Flashcard
	if err := json.Unmarshal(data, &cards); err != nil {
		return Flashcard{}, fmt.Errorf("failed to parse flashcards: %w", err)
	}

	for _, card := range cards {
		if card.ID == id {
			return card, nil
		}
	}

	return Flashcard{}, fmt.Errorf("flashcard not found: %s", id)
}

// UpdateFlashcard updates an existing flashcard
func UpdateFlashcard(card Flashcard) error {
	return SaveFlashcard(card)
}

// GetFlashcardsDueBefore returns all flashcards due before the given time
func GetFlashcardsDueBefore(time time.Time) ([]Flashcard, error) {
	filePath, err := getFilePath(cardsFile)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read flashcards file: %w", err)
	}

	var cards []Flashcard
	if err := json.Unmarshal(data, &cards); err != nil {
		return nil, fmt.Errorf("failed to parse flashcards: %w", err)
	}

	var dueCards []Flashcard
	for _, card := range cards {
		if card.NextReview.Before(time) {
			dueCards = append(dueCards, card)
		}
	}

	return dueCards, nil
}

// GetStudyStats calculates and returns study statistics
func GetStudyStats() (StudyStats, error) {
	stats := StudyStats{}

	// Get notes
	notesPath, err := getFilePath(notesFile)
	if err != nil {
		return stats, err
	}

	notesData, err := os.ReadFile(notesPath)
	if err != nil {
		return stats, fmt.Errorf("failed to read notes file: %w", err)
	}

	var notes []Note
	if err := json.Unmarshal(notesData, &notes); err != nil {
		return stats, fmt.Errorf("failed to parse notes: %w", err)
	}

	stats.TotalNotes = len(notes)

	// Get flashcards
	cardsPath, err := getFilePath(cardsFile)
	if err != nil {
		return stats, err
	}

	cardsData, err := os.ReadFile(cardsPath)
	if err != nil {
		return stats, fmt.Errorf("failed to read flashcards file: %w", err)
	}

	var cards []Flashcard
	if err := json.Unmarshal(cardsData, &cards); err != nil {
		return stats, fmt.Errorf("failed to parse flashcards: %w", err)
	}

	stats.TotalFlashcards = len(cards)

	// Count cards due today
	today := time.Now()
	tomorrow := today.Add(24 * time.Hour)

	var dueToday, learned, totalDifficulty, totalReps int

	for _, card := range cards {
		if card.NextReview.After(today) && card.NextReview.Before(tomorrow) {
			dueToday++
		}

		if card.RepCount > 0 {
			learned++
			totalDifficulty += card.Difficulty
			totalReps++
		}
	}

	stats.CardsDueToday = dueToday
	stats.CardsLearned = learned

	// Calculate review accuracy (inverted average difficulty: 0 is best, 3 is worst)
	if totalReps > 0 {
		avgDifficulty := float64(totalDifficulty) / float64(totalReps)
		stats.ReviewAccuracy = 100 * (1 - avgDifficulty/3) // Convert to percentage where 100% is best
	}

	return stats, nil
}

// GetAllFlashcards retrieves all flashcards
func GetAllFlashcards() ([]Flashcard, error) {
	filePath, err := getFilePath(cardsFile)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read flashcards file: %w", err)
	}

	var cards []Flashcard
	if err := json.Unmarshal(data, &cards); err != nil {
		return nil, fmt.Errorf("failed to parse flashcards: %w", err)
	}

	return cards, nil
}

// DeleteFlashcard removes a flashcard by ID
func DeleteFlashcard(id string) error {
	filePath, err := getFilePath(cardsFile)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read flashcards file: %w", err)
	}

	var cards []Flashcard
	if err := json.Unmarshal(data, &cards); err != nil {
		return fmt.Errorf("failed to parse flashcards: %w", err)
	}

	// Filter out the card to be deleted
	var newCards []Flashcard
	found := false
	for _, card := range cards {
		if card.ID != id {
			newCards = append(newCards, card)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("flashcard not found: %s", id)
	}

	// Save the updated list
	updatedData, err := json.MarshalIndent(newCards, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal flashcards: %w", err)
	}

	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write flashcards file: %w", err)
	}

	return nil
}

// DeleteAllFlashcards removes all flashcards
func DeleteAllFlashcards() error {
	filePath, err := getFilePath(cardsFile)
	if err != nil {
		return err
	}

	// Write an empty array to the file
	if err := os.WriteFile(filePath, []byte("[]"), 0644); err != nil {
		return fmt.Errorf("failed to clear flashcards file: %w", err)
	}

	return nil
}
