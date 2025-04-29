package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/valdezdata/md-study/internal/processor"
	"github.com/valdezdata/md-study/internal/scheduler"
	"github.com/valdezdata/md-study/internal/storage"
	"github.com/valdezdata/md-study/internal/studyengine"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "md-study",
		Short: "Study markdown files with AI-powered spaced repetition",
		Long: `A spaced repetition system that processes your markdown notes,
generates flashcards, and helps you study efficiently.`,
	}

	var importCmd = &cobra.Command{
		Use:   "import [directory]",
		Short: "Import markdown files from a directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dir := args[0]
			err := processor.ImportMarkdownFiles(dir)
			if err != nil {
				fmt.Printf("Error importing files: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully imported markdown files from %s\n", dir)
		},
	}

	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate flashcards from imported notes",
		Run: func(cmd *cobra.Command, args []string) {
			err := processor.GenerateFlashcardsForAllNotes()
			if err != nil {
				fmt.Printf("Error generating flashcards: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Successfully generated flashcards from your notes")
		},
	}

	var studyCmd = &cobra.Command{
		Use:   "study",
		Short: "Start a study session",
		Run: func(cmd *cobra.Command, args []string) {
			studyengine.StartStudySession()
		},
	}

	var statsCmd = &cobra.Command{
		Use:   "stats",
		Short: "Show study statistics",
		Run: func(cmd *cobra.Command, args []string) {
			scheduler.ShowStats()
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all flashcards",
		Run: func(cmd *cobra.Command, args []string) {
			err := processor.ListAllFlashcards()
			if err != nil {
				fmt.Printf("Error listing flashcards: %v\n", err)
				os.Exit(1)
			}
		},
	}

	var deleteCmd = &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a flashcard by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			err := storage.DeleteFlashcard(id)
			if err != nil {
				fmt.Printf("Error deleting flashcard: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Flashcard %s deleted successfully\n", id)
		},
	}

	var resetCmd = &cobra.Command{
		Use:   "reset",
		Short: "Delete all flashcards",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Are you sure you want to delete ALL flashcards? (y/n): ")
			var response string
			fmt.Scanln(&response)

			if response == "y" || response == "Y" {
				err := processor.DeleteAllFlashcards()
				if err != nil {
					fmt.Printf("Error deleting flashcards: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("All flashcards deleted successfully")
			} else {
				fmt.Println("Operation cancelled")
			}
		},
	}

	rootCmd.AddCommand(importCmd, generateCmd, studyCmd, statsCmd, listCmd, deleteCmd, resetCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
