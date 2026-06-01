package model

type GameStatus int

const (
	GamePlaying GameStatus = iota
	GameWon
	GameLost
)

// Game state
type Game struct {
	WordToGuess  string
	TriedWords   [][]LetterStatus
	CurrentGuess string
	Status       GameStatus
}

type LetterStatus struct {
	IsCorrect    bool
	IsInSolution bool
	Letter       uint8
}
