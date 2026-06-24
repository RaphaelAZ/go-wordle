package main

import (
	"log"
	"os"

	"github.com/RaphaelAZ/go-wordle/backend/internal/database"
	"github.com/RaphaelAZ/go-wordle/backend/internal/handlers"
	"github.com/RaphaelAZ/go-wordle/backend/internal/middleware"
	"github.com/RaphaelAZ/go-wordle/backend/internal/repository"
	"github.com/RaphaelAZ/go-wordle/backend/internal/seeds"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// repositories
	userRepo := repository.NewUserRepository(db)
	wordRepo := repository.NewWordRepository(db)
	gameRepo := repository.NewGameRepository(db)
	configRepo := repository.NewConfigRepository(db)

	// seed words on first run
	count, _ := wordRepo.Count()
	if count == 0 {
		n, err := wordRepo.Seed(seeds.FrenchWords)
		if err != nil {
			log.Printf("word seed warning: %v", err)
		} else {
			log.Printf("seeded %d words", n)
		}
	}

	// seed users on first run
	userCount, _ := userRepo.Count()
	if userCount == 0 {
		n, err := seeds.SeedUsers(userRepo, seeds.Users)
		if err != nil {
			log.Printf("user seed warning: %v", err)
		} else {
			log.Printf("seeded %d users", n)
		}
	}

	// handlers
	authHandler := handlers.NewAuthHandler(userRepo)
	wordHandler := handlers.NewWordHandler(wordRepo)
	gameHandler := handlers.NewGameHandler(gameRepo)
	configHandler := handlers.NewConfigHandler(configRepo)

	r := gin.Default()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/me", authHandler.Me)

			protected.GET("/words/random", wordHandler.Random)

			protected.POST("/games", gameHandler.Create)
			protected.GET("/games", gameHandler.List)
			protected.GET("/games/stats", gameHandler.Stats)

			protected.GET("/config", configHandler.Get)
			protected.PUT("/config", configHandler.Upsert)
		}
	}

	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ip + ":" + port
	log.Printf("server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
