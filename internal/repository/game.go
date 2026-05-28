package repository

import (
	"database/sql"
	"fmt"

	"github.com/RaphaelAZ/go-wordle/internal/models"
)

type GameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) Create(userID, wordID int, attempts []byte, won bool, duration int) (*models.GameSession, error) {
	g := &models.GameSession{}
	err := r.db.QueryRow(
		`INSERT INTO game_sessions (user_id, word_id, attempts, won, duration)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, word_id, attempts, won, duration, created_at`,
		userID, wordID, attempts, won, duration,
	).Scan(&g.ID, &g.UserID, &g.WordID, &g.Attempts, &g.Won, &g.Duration, &g.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create game session: %w", err)
	}
	return g, nil
}

func (r *GameRepository) ListByUser(userID int) ([]models.GameSession, error) {
	rows, err := r.db.Query(
		`SELECT gs.id, gs.user_id, gs.word_id, w.word, gs.attempts, gs.won, gs.duration, gs.created_at
		 FROM game_sessions gs
		 JOIN words w ON w.id = gs.word_id
		 WHERE gs.user_id = $1
		 ORDER BY gs.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list games: %w", err)
	}
	defer rows.Close()

	var games []models.GameSession
	for rows.Next() {
		var g models.GameSession
		if err := rows.Scan(&g.ID, &g.UserID, &g.WordID, &g.Word, &g.Attempts, &g.Won, &g.Duration, &g.CreatedAt); err != nil {
			return nil, err
		}
		games = append(games, g)
	}
	return games, rows.Err()
}

func (r *GameRepository) Stats(userID int) (*models.GameStats, error) {
	stats := &models.GameStats{}
	err := r.db.QueryRow(
		`SELECT
			COUNT(*)                                          AS total_games,
			COUNT(*) FILTER (WHERE won)                      AS wins,
			COALESCE(AVG(duration), 0)                       AS avg_duration
		 FROM game_sessions
		 WHERE user_id = $1`,
		userID,
	).Scan(&stats.TotalGames, &stats.Wins, &stats.AvgDuration)
	if err != nil {
		return nil, fmt.Errorf("get stats: %w", err)
	}

	if stats.TotalGames > 0 {
		stats.WinRate = float64(stats.Wins) / float64(stats.TotalGames) * 100
	}

	stats.CurrentStreak, stats.BestStreak = r.computeStreaks(userID)

	return stats, nil
}

func (r *GameRepository) computeStreaks(userID int) (current, best int) {
	rows, err := r.db.Query(
		`SELECT won FROM game_sessions WHERE user_id = $1 ORDER BY created_at ASC`,
		userID,
	)
	if err != nil {
		return 0, 0
	}
	defer rows.Close()

	var results []bool
	for rows.Next() {
		var won bool
		rows.Scan(&won)
		results = append(results, won)
	}

	// best streak (chronological)
	streak := 0
	for _, won := range results {
		if won {
			streak++
			if streak > best {
				best = streak
			}
		} else {
			streak = 0
		}
	}

	// current streak (from most recent game going back)
	for i := len(results) - 1; i >= 0; i-- {
		if results[i] {
			current++
		} else {
			break
		}
	}

	return current, best
}
