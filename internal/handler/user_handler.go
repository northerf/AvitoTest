package handler

import (
	"Avito/internal/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) SetUserActive(c *gin.Context) {
	var req schema.UpdateUserActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "Invalid request body")
		return
	}
	userID := c.Query("user_id")
	if userID == "" {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "user_id is required")
		return
	}
	err := h.services.User.SetActive(c.Request.Context(), userID, *req.IsActive)
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	user, err := h.services.User.GetByID(c.Request.Context(), userID)
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, schema.User{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	})
}

func (h *Handler) GetUserReviews(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "user_id is required")
		return
	}

	prs, err := h.services.PullRequest.GetByReviewerID(c.Request.Context(), userID)
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	result := make([]schema.PullRequestShort, len(prs))
	for i, pr := range prs {
		result[i] = schema.PullRequestShort{
			PullRequestID:   pr.ID,
			PullRequestName: pr.Name,
			AuthorID:        pr.AuthorID,
			Status:          schema.PRStatus(pr.Status),
		}
	}
	c.JSON(http.StatusOK, result)
}
