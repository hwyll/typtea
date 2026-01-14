package cmd

import (
	"fmt"
	"strings"

	"github.com/hwyll/typtea/internal/game"
	"github.com/hwyll/typtea/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	duration  int    // Duration of the typing test in seconds
	language  string // Language for the typing test, default is "en"
	listLangs bool   // Flag to list all available languages
	wordCount int    // Number of words for word mode
	mode      string // Game mode: "time" or "word"
)

// startCmd represents the start command for the typing test
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a typing test",
	Long:  "Start a new typing test session with customizable duration, word count, and language",
	Example: `  # Time mode (default)
  typtea start --duration 60 --lang python
  typtea start -d 30 -l javascript

  # Word mode
  typtea start --mode word --words 100
  typtea start -m word -w 50 -l go

  # List available languages
  typtea start --list-langs`,
	RunE: runTypingTest,
}

func init() {
	startCmd.Flags().IntVarP(&duration, "duration", "d", 30, "Test duration in seconds (10-300) for time mode")
	startCmd.Flags().StringVarP(&language, "lang", "l", "en", "Language for typing test")
	startCmd.Flags().BoolVar(&listLangs, "list-langs", false, "List all available languages")
	startCmd.Flags().IntVarP(&wordCount, "words", "w", 50, "Number of words (10-500) for word mode")
	startCmd.Flags().StringVarP(&mode, "mode", "m", "time", "Game mode: 'time' or 'word'")
}

// runTypingTest runs the typing test or lists languages if requested
func runTypingTest(cmd *cobra.Command, args []string) error {
	// Initialize the language manager
	langManager := game.NewLanguageManager()

	// If --list-langs flag is set, print available languages and exit
	if listLangs {
		cmd.Println("Available languages:")
		for _, lang := range langManager.GetAvailableLanguages() {
			cmd.Printf("  %s\n", lang)
		}
		return nil
	}

	// Validate language availability
	if !langManager.IsLanguageAvailable(language) {
		available := langManager.GetAvailableLanguages()
		cmd.PrintErrf("Error: Language '%s' not available.\n", language)
		cmd.PrintErrf("Available languages: %s\n", strings.Join(available, ", "))
		return fmt.Errorf("invalid language: %s", language)
	}

	// Validate and determine game mode
	var gameMode game.GameMode
	switch mode {
	case "time":
		if duration < 10 || duration > 300 {
			return fmt.Errorf("duration must be between 10 and 300 seconds")
		}
		gameMode = game.Time
	case "word":
		if wordCount < 10 || wordCount > 500 {
			return fmt.Errorf("word count must be between 10 and 500")
		}
		gameMode = game.Word
	default:
		return fmt.Errorf("invalid mode: %s (must be 'time' or 'word')", mode)
	}

	// Create a new typing test model
	model, err := tui.NewModel(gameMode, duration, wordCount, language)
	if err != nil {
		return fmt.Errorf("error creating typing test: %w", err)
	}

	// Start the TUI program with alternate screen
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI program: %w", err)
	}

	return nil
}
