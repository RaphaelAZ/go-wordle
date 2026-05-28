package models

import (
	"encoding/json"
	"time"
)

type GameSession struct {
	ID        int             `json:"id"`
	UserID    int             `json:"user_id"`
	WordID    int             `json:"word_id"`
	Word      string          `json:"word,omitempty"`
	Attempts  json.RawMessage `json:"attempts"`
	Won       bool            `json:"won"`
	Duration  int             `json:"duration"`
	CreatedAt time.Time       `json:"created_at"`
}

type GameStats struct {
	TotalGames  int     `json:"total_games"`
	Wins        int     `json:"wins"`
	WinRate     float64 `json:"win_rate"`
	CurrentStreak int   `json:"current_streak"`
	BestStreak  int     `json:"best_streak"`
	AvgDuration float64 `json:"avg_duration"`
}
