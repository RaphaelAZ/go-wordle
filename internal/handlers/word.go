package handlers

import (
	"net/http"

	"github.com/RaphaelAZ/go-wordle/internal/repository"
	"github.com/gin-gonic/gin"
)

type WordHandler struct {
	words *repository.WordRepository
}

func NewWordHandler(words *repository.WordRepository) *WordHandler {
	return &WordHandler{words: words}
}

func (h *WordHandler) Random(c *gin.Context) {
	word, err := h.words.Random()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, word)
}
