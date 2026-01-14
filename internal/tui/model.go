package tui

import (
	"fmt"
	"time"

	"github.com/hwyll/typtea/internal/game"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the state of the typing test application
type Model struct {
	game        *game.TypingGame
	width       int
	height      int
	showResults bool
	finalStats  game.TypingStats
	duration    int
	wordCount   int
	language    string
}

// tickMsg is a message type used to handle periodic updates in the application
type tickMsg time.Time

// NewModel initializes a new Model instance with the specified mode, duration, word count, and language
func NewModel(mode game.GameMode, duration int, wordCount int, language string) (*Model, error) {
	if err := game.SetLanguage(language); err != nil {
		return nil, fmt.Errorf("failed to load language '%s': %v", language, err)
	}

	var typingGame *game.TypingGame
	switch mode {
	case game.Time:
		typingGame = game.NewTimeGame(duration)
	case game.Word:
		typingGame = game.NewWordGame(wordCount)
	default:
		return nil, fmt.Errorf("invalid game mode")
	}

	return &Model{
		game:      typingGame,
		duration:  duration,
		wordCount: wordCount,
		language:  language,
	}, nil
}

// restartTest resets the game state for a new typing test session
func (m *Model) restartTest() {
	switch m.game.Mode {
	case game.Time:
		m.game = game.NewTimeGame(m.duration)
	case game.Word:
		m.game = game.NewWordGame(m.wordCount)
	}
	m.showResults = false
	m.finalStats = game.TypingStats{}
}

// Init initializes the model and starts the tick command for periodic updates
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		tea.EnterAltScreen,
	)
}

// tickCmd returns a command that sends a tick message every 100 milliseconds
func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
