package model

import "time"

type GameStatus int

const (
	GamePlaying GameStatus = iota
	GameWon
	GameLost
)

// Game state
type Game struct {
	WordID       int
	WordToGuess  string
	TriedWords   [][]LetterStatus
	CurrentGuess string
	Status       GameStatus
	UsedLetters  map[byte]LetterStatus
	WordLoading  bool
	WordError    string
	StartedAt    time.Time
	SaveLoading  bool
	SaveError    string
}

type LetterStatus struct {
	IsCorrect    bool
	IsInSolution bool
	Letter       uint8
}

type RestartGameMsg struct{}
