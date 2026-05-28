package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RaphaelAZ/go-wordle/internal/repository"
	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	configs *repository.ConfigRepository
}

func NewConfigHandler(configs *repository.ConfigRepository) *ConfigHandler {
	return &ConfigHandler{configs: configs}
}

type upsertConfigRequest struct {
	Config json.RawMessage `json:"config" binding:"required"`
	State  json.RawMessage `json:"state"  binding:"required"`
}

func (h *ConfigHandler) Get(c *gin.Context) {
	userID := c.GetInt("user_id")
	cfg, err := h.configs.Get(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cfg == nil {
		c.JSON(http.StatusOK, gin.H{"config": json.RawMessage("{}"), "state": json.RawMessage("{}")})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func (h *ConfigHandler) Upsert(c *gin.Context) {
	var req upsertConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")
	cfg, err := h.configs.Upsert(userID, req.Config, req.State)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}
