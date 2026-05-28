package models

import "time"

type Word struct {
	ID        int       `json:"id"`
	Word      string    `json:"word"`
	CreatedAt time.Time `json:"created_at"`
}
