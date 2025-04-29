# MD-Study: AI-Powered Spaced Repetition for Markdown Notes

MD-Study is a command-line tool that automatically converts your Markdown notes into flashcards for spaced repetition study. It leverages AI to extract key concepts from your notes and transforms them into effective question-answer pairs.

## Features

- Import existing markdown files as study material
- AI-powered flashcard generation from your notes
- Spaced repetition algorithm for efficient learning
- Interactive CLI-based study sessions
- Track learning progress with statistics
- Manage your flashcards (list, delete, reset)

## Installation

### Prerequisites

- Go 1.18 or later
- OpenAI API key (export as OPENAI_API_KEY)

### Building from source

```bash
# Clone the repository
git clone https://github.com/yourusername/md-study.git
cd md-study

# Build the application
make build

# Or install directly
make install
```

## Usage

### Workflow

1. **Import** markdown files
2. **Generate** flashcards using AI (will only process notes without existing flashcards)
3. **Study** with spaced repetition
4. View study **stats**

### Commands

```bash
# Import your markdown files
md-study import /path/to/markdown/files

# Generate flashcards using AI
md-study generate

# Start a study session
md-study study

# View your study statistics
md-study stats

# List all flashcards
md-study list

# Delete a specific flashcard
md-study delete [flashcard-id]

# Reset all flashcards
md-study reset
```

### Environment Variables

- `OPENAI_API_KEY`: Your OpenAI API key (required)

## How It Works

MD-Study uses the SM-2 spaced repetition algorithm to schedule reviews based on your performance. When you rate a flashcard during study, the system adjusts the next review time accordingly:

- **Easy**: Longer intervals between reviews
- **Good**: Standard intervals
- **Hard**: Shorter intervals
- **Again**: Very short intervals for immediate reinforcement

The AI-powered flashcard generation analyzes your markdown notes to identify key concepts and creates question-answer pairs that effectively test your understanding of the material.

## Data Storage

All your notes, flashcards, and study statistics are stored in `~/.md-study/` directory:

- `notes.json`: Imported markdown files
- `flashcards.json`: Generated flashcards with spaced repetition metadata
- `stats.json`: Study progress and statistics

## Future Improvements

- Web UI for more interactive study
- Support for images and formatting in flashcards
- Custom tagging of notes for targeted study
- Improved markdown parsing to extract headings, lists, etc.
- Multiple AI providers

## License

MIT
