package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/RaphaelAZ/go-wordle/internal/models"
)

type WordRepository struct {
	db *sql.DB
}

func NewWordRepository(db *sql.DB) *WordRepository {
	return &WordRepository{db: db}
}

func (r *WordRepository) Random() (*models.Word, error) {
	word := &models.Word{}
	err := r.db.QueryRow(
		`SELECT id, word, created_at FROM words ORDER BY RANDOM() LIMIT 1`,
	).Scan(&word.ID, &word.Word, &word.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no words in database")
	}
	if err != nil {
		return nil, fmt.Errorf("get random word: %w", err)
	}
	return word, nil
}

func (r *WordRepository) Seed(words []string) (int, error) {
	inserted := 0
	for _, w := range words {
		res, err := r.db.Exec(
			`INSERT INTO words (word) VALUES ($1) ON CONFLICT (word) DO NOTHING`,
			w,
		)
		if err != nil {
			return inserted, fmt.Errorf("seed word %q: %w", w, err)
		}
		n, _ := res.RowsAffected()
		inserted += int(n)
	}
	return inserted, nil
}

func (r *WordRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM words`).Scan(&count)
	return count, err
}
