package processor

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/valdezdata/md-study/internal/storage"
)

// GenerateFlashcards uses AI to create flashcards from notes
func GenerateFlashcards(noteID string) ([]storage.Flashcard, error) {
	note, err := storage.GetNote(noteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(apiKey)

	// Construct the prompt
	prompt := fmt.Sprintf("Create 5 flashcards in question-answer format from the following notes. Format each as 'Q: [question]' on one line and 'A: [answer]' on another line.\n\nNotes:\n%s", note.RawContent)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4.1-nano", // You can change to your preferred model
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "system",
					Content: "You are a helpful assistant that creates effective flashcards for learning.",
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
			Temperature: 0.3, // Lower temperature for more consistent output
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	// Parse the AI response into flashcards
	return parseFlashcardsFromResponse(resp.Choices[0].Message.Content, noteID)
}

// parseFlashcardsFromResponse extracts Q&A pairs from the AI response
func parseFlashcardsFromResponse(response, noteID string) ([]storage.Flashcard, error) {
	lines := strings.Split(response, "\n")
	var flashcards []storage.Flashcard

	var currentQuestion, currentAnswer string
	for i, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Q:") {
			// If we already have a question and answer, save the flashcard
			if currentQuestion != "" && currentAnswer != "" {
				flashcards = append(flashcards, storage.Flashcard{
					NoteID:     noteID,
					Question:   currentQuestion,
					Answer:     currentAnswer,
					Difficulty: 0, // Initial difficulty
					NextReview: time.Now(),
				})
				currentQuestion, currentAnswer = "", ""
			}

			currentQuestion = strings.TrimPrefix(line, "Q:")
			currentQuestion = strings.TrimSpace(currentQuestion)
		} else if strings.HasPrefix(line, "A:") {
			currentAnswer = strings.TrimPrefix(line, "A:")
			currentAnswer = strings.TrimSpace(currentAnswer)

			// If this is the last line or the next line will be a new question, save the flashcard
			if (i == len(lines)-1 || strings.HasPrefix(strings.TrimSpace(lines[i+1]), "Q:")) && currentQuestion != "" {
				flashcards = append(flashcards, storage.Flashcard{
					NoteID:     noteID,
					Question:   currentQuestion,
					Answer:     currentAnswer,
					Difficulty: 0, // Initial difficulty
					NextReview: time.Now(),
				})
				currentQuestion, currentAnswer = "", ""
			}
		}
	}

	return flashcards, nil
}

// GenerateFlashcardsForAllNotes processes all imported notes and creates flashcards
func GenerateFlashcardsForAllNotes() error {
	// Get all notes
	notes, err := storage.GetAllNotes()
	if err != nil {
		return fmt.Errorf("failed to get notes: %w", err)
	}

	if len(notes) == 0 {
		return fmt.Errorf("no notes found - please import some markdown files first")
	}

	fmt.Printf("Found %d notes to process\n", len(notes))

	// First, get all existing flashcards
	existingCards, err := storage.GetAllFlashcards()
	if err != nil {
		return fmt.Errorf("failed to get existing flashcards: %w", err)
	}

	// Create a map to track which notes already have flashcards
	notesWithCards := make(map[string]bool)
	for _, card := range existingCards {
		notesWithCards[card.NoteID] = true
	}

	// Track how many notes were processed
	processedCount := 0

	for i, note := range notes {
		// Skip notes that already have flashcards
		if notesWithCards[note.ID] {
			fmt.Printf("[%d/%d] Skipping %s (already has flashcards)\n", i+1, len(notes), note.Filename)
			continue
		}

		fmt.Printf("[%d/%d] Generating flashcards for %s...\n", i+1, len(notes), note.Filename)
		processedCount++

		// Generate flashcards for this note
		flashcards, err := GenerateFlashcards(note.ID)
		if err != nil {
			return fmt.Errorf("failed to generate flashcards for %s: %w", note.Filename, err)
		}

		// Save each flashcard
		for _, card := range flashcards {
			if err := storage.SaveFlashcard(card); err != nil {
				return fmt.Errorf("failed to save flashcard: %w", err)
			}
		}

		fmt.Printf("  Created %d flashcards\n", len(flashcards))
	}

	if processedCount == 0 {
		fmt.Println("No new notes to process. All notes already have flashcards.")
	}

	return nil
}

// ListAllFlashcards displays all flashcards in the system
func ListAllFlashcards() error {
	// Get all flashcards
	cards, err := storage.GetAllFlashcards()
	if err != nil {
		return fmt.Errorf("failed to get flashcards: %w", err)
	}

	if len(cards) == 0 {
		fmt.Println("No flashcards found. Use the 'generate' command to create some.")
		return nil
	}

	fmt.Printf("Found %d flashcards:\n\n", len(cards))

	for i, card := range cards {
		fmt.Printf("Flashcard #%d:\n", i+1)
		fmt.Printf("Question: %s\n", card.Question)
		fmt.Printf("Answer: %s\n", card.Answer)
		fmt.Printf("Next review: %s\n", card.NextReview.Format("2006-01-02 15:04:05"))
		fmt.Println("---------------------------------------")
	}

	return nil
}

// DeleteAllFlashcards removes all flashcards
func DeleteAllFlashcards() error {
	return storage.DeleteAllFlashcards()
}
