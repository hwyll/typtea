package game

import (
	"strings"
	"time"
)

// GameMode represents the type of typing test
type GameMode int

const (
	Time GameMode = iota // Time-based test
	Word                 // Word count-based test
)

func (m GameMode) String() string {
	switch m {
	case Time:
		return "time"
	case Word:
		return "word"
	default:
		return "unknown"
	}
}

// TypingStats holds the statistics for a game session
type TypingStats struct {
	WPM               float64
	Accuracy          float64
	CharactersTyped   int
	CorrectChars      int
	TotalChars        int
	TimeElapsed       time.Duration
	IsComplete        bool
	UncorrectedErrors int
}

// TypingGame represents the state of a game session
type TypingGame struct {
	AllWords        []string
	DisplayLines    []string
	UserInput       string
	CurrentPos      int
	GlobalPos       int
	StartTime       time.Time
	IsStarted       bool
	Errors          map[int]bool
	TotalErrorsMade int
	LinesPerView    int
	CharsPerLine    int
	WordsTyped      int

	// Mode configuration
	Mode        GameMode
	Duration    int // For Time mode (seconds)
	TargetWords int // For Word mode (word count)
}

// NewTimeGame creates a time-based typing game
func NewTimeGame(duration int) *TypingGame {
	game := &TypingGame{
		AllWords:     GenerateWords(200),
		Duration:     duration,
		Mode:         Time,
		Errors:       make(map[int]bool),
		LinesPerView: 3,
		CharsPerLine: 50,
	}
	game.generateDisplayLines()
	return game
}

// NewWordGame creates a word count-based typing game
func NewWordGame(wordCount int) *TypingGame {
	game := &TypingGame{
		AllWords:     GenerateWords(wordCount),
		TargetWords:  wordCount,
		Mode:         Word,
		Errors:       make(map[int]bool),
		LinesPerView: 3,
		CharsPerLine: 50,
	}
	game.generateDisplayLines()
	return game
}

// Reset reinitializes the game to a fresh state
func (g *TypingGame) Reset() {
	switch g.Mode {
	case Time:
		*g = *NewTimeGame(g.Duration)
	case Word:
		*g = *NewWordGame(g.TargetWords)
	}
}

// generateDisplayLines creates the initial display lines based on the words available
func (g *TypingGame) generateDisplayLines() {
	lines := make([]string, 0, g.LinesPerView)
	wordIndex := g.WordsTyped

	// Generate exactly g.LinesPerView lines
	for lineNum := 0; lineNum < g.LinesPerView && wordIndex < len(g.AllWords); lineNum++ {
		var currentLine strings.Builder

		// Fill current line with words
		for wordIndex < len(g.AllWords) {
			word := g.AllWords[wordIndex]
			spaceNeeded := 0
			if currentLine.Len() > 0 {
				spaceNeeded = 1
			}

			// Check if word fits
			if currentLine.Len()+spaceNeeded+len(word) <= g.CharsPerLine {
				if currentLine.Len() > 0 {
					currentLine.WriteString(" ")
				}
				currentLine.WriteString(word)
				wordIndex++
			} else {
				// Word doesn't fit, break to next line
				break
			}
		}

		// Add the completed line
		if currentLine.Len() > 0 {
			lines = append(lines, currentLine.String())
		} else {
			// If no words fit, add empty line
			lines = append(lines, "")
		}
	}

	// Ensure we have exactly g.LinesPerView lines
	for len(lines) < g.LinesPerView {
		lines = append(lines, "")
	}

	// Truncate if somehow we have more than g.LinesPerView lines
	if len(lines) > g.LinesPerView {
		lines = lines[:g.LinesPerView]
	}

	g.DisplayLines = lines
}

// Start initializes the game session if it hasn't started yet
func (g *TypingGame) Start() {
	if !g.IsStarted {
		g.StartTime = time.Now()
		g.IsStarted = true
	}
}

// IsFinished checks if the game has completed based on the current mode
func (g *TypingGame) IsFinished() bool {
	if !g.IsStarted {
		return false
	}

	switch g.Mode {
	case Time:
		return g.IsTimeUp()
	case Word:
		return g.GetWordsTyped() >= g.TargetWords
	default:
		return false
	}
}

// AddCharacter handles user input and updates game state
func (g *TypingGame) AddCharacter(char rune) {
	if !g.IsStarted {
		g.Start()
	}

	if g.IsFinished() {
		return
	}

	lineText := []rune(g.DisplayLines[0])

	// If at end of line, only shift if user just typed space
	if g.CurrentPos == len(lineText) {
		if char == ' ' {
			g.UserInput += string(char)
			g.CurrentPos++
			g.GlobalPos++
			g.shiftLines()
		}
		return
	}

	// Normal character processing
	if g.CurrentPos < len(lineText) && g.CurrentPos >= 0 {
		g.UserInput += string(char)
		if lineText[g.CurrentPos] != char {
			g.Errors[g.GlobalPos] = true
			g.TotalErrorsMade++
		}
		g.CurrentPos++
		g.GlobalPos++
	}
}

// HandleEnterKey handles Enter key press for line progression
func (g *TypingGame) HandleEnterKey() bool {
	if g.IsFinished() {
		return false
	}

	lineText := []rune(g.DisplayLines[0])

	// Only allow Enter to progress if at end of line
	if g.CurrentPos == len(lineText) {
		// Treat Enter like Space internally for consistency
		g.UserInput += " "
		g.CurrentPos++
		g.GlobalPos++
		g.shiftLines()
		return true
	}

	return false
}

// shiftLines moves to the next line in the game, updating the words typed and generating new lines
func (g *TypingGame) shiftLines() {
	// Move to next line
	g.WordsTyped += len(strings.Fields(g.DisplayLines[0]))
	g.CurrentPos = 0

	// Generate new lines
	g.generateDisplayLines()

	// Only extend words in Time mode (infinite words)
	if g.Mode == Time && g.WordsTyped > len(g.AllWords)-50 {
		newWords := GenerateWords(100)
		g.AllWords = append(g.AllWords, newWords...)
	}
}

// RemoveCharacter removes the last character from the user input and updates the position
func (g *TypingGame) RemoveCharacter() {
	if len(g.UserInput) > 0 && g.CurrentPos > 0 {
		g.UserInput = g.UserInput[:len(g.UserInput)-1]
		g.CurrentPos--
		g.GlobalPos--

		// Remove error mark if previously added
		delete(g.Errors, g.GlobalPos)
	}
}

// GetDisplayText returns the current text to be displayed in the game
func (g *TypingGame) GetDisplayText() string {
	return strings.Join(g.DisplayLines, " ")
}

// GetWordsTyped returns the current number of words typed (only counts complete words)
func (g *TypingGame) GetWordsTyped() int {
	// Words from completed lines
	completedWords := g.WordsTyped

	// Count complete words on the current line (words followed by space)
	if g.CurrentPos > 0 && len(g.DisplayLines) > 0 {
		currentLineText := g.DisplayLines[0]
		if g.CurrentPos <= len(currentLineText) {
			typedPortion := currentLineText[:g.CurrentPos]

			// Count words only if the last character is a space
			// This means the word is complete
			trimmed := strings.TrimRight(typedPortion, " ")
			if len(trimmed) > 0 {
				words := strings.Fields(trimmed)
				// Only count if we've typed a space after the last word
				if typedPortion != trimmed {
					return completedWords + len(words)
				}
				// Still typing the last word, don't count it
				if len(words) > 0 {
					return completedWords + len(words) - 1
				}
			}
		}
	}

	return completedWords
}

// GetStats calculates and returns the typing statistics for the current game session
func (g *TypingGame) GetStats() TypingStats {
	if !g.IsStarted {
		return TypingStats{}
	}

	elapsed := time.Since(g.StartTime)
	minutes := elapsed.Minutes()

	// Calculate Gross WPM (all typed entries / 5 / time in minutes)
	grossWPM := 0.0
	if minutes > 0 {
		grossWPM = float64(g.GlobalPos) / 5 / minutes
	}

	// Calculate uncorrected errors (errors still present in the text)
	uncorrectedErrors := len(g.Errors)

	// Calculate Net WPM (Gross WPM - uncorrected errors per minute)
	errorRate := float64(uncorrectedErrors) / minutes
	netWPM := grossWPM - errorRate

	// Ensure Net WPM doesn't go below 0
	if netWPM < 0 {
		netWPM = 0
	}

	// Calculate accuracy (correct characters / total characters typed * 100)
	correctChars := g.GlobalPos - g.TotalErrorsMade
	accuracy := 0.0
	if g.GlobalPos > 0 {
		accuracy = float64(correctChars) / float64(g.GlobalPos) * 100
	}

	// Ensure accuracy doesn't go below 0
	if accuracy < 0 {
		accuracy = 0
	}

	return TypingStats{
		WPM:               netWPM,
		Accuracy:          accuracy,
		CharactersTyped:   g.GlobalPos,
		CorrectChars:      correctChars,
		TotalChars:        len([]rune(g.GetDisplayText())),
		TimeElapsed:       elapsed,
		IsComplete:        g.IsFinished(),
		UncorrectedErrors: uncorrectedErrors,
	}
}

// IsTimeUp checks if the game time has exceeded the specified duration (Time mode only)
func (g *TypingGame) IsTimeUp() bool {
	if !g.IsStarted || g.Mode != Time {
		return false
	}
	return time.Since(g.StartTime).Seconds() >= float64(g.Duration)
}

// GetRemainingTime returns the remaining time in seconds for Time mode
func (g *TypingGame) GetRemainingTime() int {
	if g.Mode != Time {
		return 0
	}

	if !g.IsStarted {
		return g.Duration
	}

	elapsed := int(time.Since(g.StartTime).Seconds())
	remaining := g.Duration - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetRemainingWords returns the remaining words for Word mode
func (g *TypingGame) GetRemainingWords() int {
	if g.Mode != Word {
		return 0
	}

	remaining := g.TargetWords - g.GetWordsTyped()
	if remaining < 0 {
		return 0
	}
	return remaining
}
