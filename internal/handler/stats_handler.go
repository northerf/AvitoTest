package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Avito/internal/schema"
)

func (h *Handler) GetStatistics(c *gin.Context) {
	stats, err := h.services.Stats.GetStatistics(c.Request.Context())
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetUserStatistics(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "user_id is required")
		return
	}
	stats, err := h.services.Stats.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, stats)
}
