package scheduler

import (
	"fmt"
	"time"

	"github.com/valdezdata/md-study/internal/storage"
)

// SM-2 algorithm intervals (in hours)
var intervals = [][]int{
	{0, 24, 144, 432}, // Easy intervals: 0h, 1d, 6d, 18d
	{0, 8, 48, 172},   // Good intervals: 0h, 8h, 2d, 7d
	{0, 3, 24, 72},    // Hard intervals: 0h, 3h, 1d, 3d
	{0, 1, 3, 8},      // Again intervals: 0h, 1h, 3h, 8h
}

// GetDueFlashcards returns flashcards due for review
func GetDueFlashcards() ([]storage.Flashcard, error) {
	return storage.GetFlashcardsDueBefore(time.Now())
}

// UpdateFlashcard updates a flashcard's difficulty and next review time
func UpdateFlashcard(id string, difficulty int) error {
	card, err := storage.GetFlashcard(id)
	if err != nil {
		return err
	}

	// Calculate next review time based on difficulty and repetition count
	nextInterval := intervals[difficulty][min(card.RepCount, len(intervals[difficulty])-1)]
	card.NextReview = time.Now().Add(time.Duration(nextInterval) * time.Hour)
	card.RepCount++
	card.Difficulty = difficulty

	return storage.UpdateFlashcard(card)
}

// ShowStats displays study statistics
func ShowStats() {
	stats, err := storage.GetStudyStats()
	if err != nil {
		fmt.Printf("Error getting stats: %v\n", err)
		return
	}

	fmt.Println("Study Statistics:")
	fmt.Printf("Total notes: %d\n", stats.TotalNotes)
	fmt.Printf("Total flashcards: %d\n", stats.TotalFlashcards)
	fmt.Printf("Cards due today: %d\n", stats.CardsDueToday)
	fmt.Printf("Cards learned: %d\n", stats.CardsLearned)
	fmt.Printf("Review accuracy: %.1f%%\n", stats.ReviewAccuracy)
}

// min returns the smaller of a and b
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
