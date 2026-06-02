package seeds

import (
	"fmt"

	"github.com/RaphaelAZ/go-wordle/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type SeedUser struct {
	Username string
	Email    string
	Password string
}

var Users = []SeedUser{
	{Username: "alice", Email: "alice@example.com", Password: "password123"},
	{Username: "bob", Email: "bob@example.com", Password: "password123"},
	{Username: "charlie", Email: "charlie@example.com", Password: "password123"},
}

func SeedUsers(repo *repository.UserRepository, users []SeedUser) (int, error) {
	seeded := 0
	for _, u := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return seeded, fmt.Errorf("hash password for %s: %w", u.Username, err)
		}
		if _, err := repo.Create(u.Username, u.Email, string(hash)); err != nil {
			return seeded, fmt.Errorf("seed user %s: %w", u.Username, err)
		}
		seeded++
	}
	return seeded, nil
}
