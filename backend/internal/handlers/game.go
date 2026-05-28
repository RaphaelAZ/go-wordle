package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RaphaelAZ/go-wordle/backend/internal/models"
	"github.com/RaphaelAZ/go-wordle/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	games *repository.GameRepository
}

func NewGameHandler(games *repository.GameRepository) *GameHandler {
	return &GameHandler{games: games}
}

type createGameRequest struct {
	WordID   int             `json:"word_id"  binding:"required"`
	Attempts json.RawMessage `json:"attempts" binding:"required"`
	Won      bool            `json:"won"`
	Duration int             `json:"duration"`
}

func (h *GameHandler) Create(c *gin.Context) {
	var req createGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")
	game, err := h.games.Create(userID, req.WordID, req.Attempts, req.Won, req.Duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, game)
}

func (h *GameHandler) List(c *gin.Context) {
	userID := c.GetInt("user_id")
	games, err := h.games.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if games == nil {
		games = []models.GameSession{}
	}
	c.JSON(http.StatusOK, games)
}

func (h *GameHandler) Stats(c *gin.Context) {
	userID := c.GetInt("user_id")
	stats, err := h.games.Stats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
