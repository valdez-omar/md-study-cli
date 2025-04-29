package studyengine

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/valdezdata/md-study/internal/scheduler"
)

// StartStudySession begins an interactive study session
func StartStudySession() {
	// Get due flashcards
	flashcards, err := scheduler.GetDueFlashcards()
	if err != nil {
		fmt.Printf("Error getting flashcards: %v\n", err)
		return
	}

	if len(flashcards) == 0 {
		fmt.Println("No flashcards due for review right now!")
		return
	}

	fmt.Printf("Starting study session with %d flashcards\n", len(flashcards))

	reader := bufio.NewReader(os.Stdin)

	for i, card := range flashcards {
		fmt.Printf("\n--- Card %d/%d ---\n", i+1, len(flashcards))
		color.Cyan("%s", card.Question)

		fmt.Print("\nPress Enter to see answer...")
		reader.ReadString('\n')

		color.Yellow("%s", card.Answer)

		fmt.Println("\nRate your recall:")
		color.Green("1 - Easy")
		color.Cyan("2 - Good")
		color.Yellow("3 - Hard")
		color.Red("4 - Again")

		fmt.Print("\nYour rating (1-4): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var difficulty int
		switch input {
		case "1":
			difficulty = 0 // Easy
		case "2":
			difficulty = 1 // Good
		case "3":
			difficulty = 2 // Hard
		case "4":
			difficulty = 3 // Again
		default:
			difficulty = 2 // Default to Hard if invalid input
		}

		// Update card difficulty and next review time
		scheduler.UpdateFlashcard(card.ID, difficulty)
	}

	fmt.Println("\nStudy session complete!")
}
