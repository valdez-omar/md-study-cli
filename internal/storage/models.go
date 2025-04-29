package storage

import (
	"time"
)

// Note represents a markdown note
type Note struct {
	ID         string    `json:"id"`
	FilePath   string    `json:"file_path"`
	Filename   string    `json:"filename"`
	RawContent string    `json:"raw_content"`
	LastImport time.Time `json:"last_import"`
	Flashcards []string  `json:"flashcard_ids"`
}

// Flashcard represents a question-answer pair for studying
type Flashcard struct {
	ID         string    `json:"id"`
	NoteID     string    `json:"note_id"`
	Question   string    `json:"question"`
	Answer     string    `json:"answer"`
	Difficulty int       `json:"difficulty"` // 0-3: Easy, Good, Hard, Again
	RepCount   int       `json:"rep_count"`  // Number of repetitions
	LastReview time.Time `json:"last_review"`
	NextReview time.Time `json:"next_review"`
}

// StudyStats represents study statistics
type StudyStats struct {
	TotalNotes      int     `json:"total_notes"`
	TotalFlashcards int     `json:"total_flashcards"`
	CardsDueToday   int     `json:"cards_due_today"`
	CardsLearned    int     `json:"cards_learned"`
	ReviewAccuracy  float64 `json:"review_accuracy"`
}
