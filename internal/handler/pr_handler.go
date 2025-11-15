package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Avito/internal/schema"
)

func (h *Handler) CreatePR(c *gin.Context) {
	var req schema.CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "Invalid request body")
		return
	}
	pr, err := h.services.PullRequest.Create(c.Request.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	response := schema.PullRequest{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            schema.PRStatus(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
	c.JSON(http.StatusCreated, gin.H{
		"pr": response,
	})
}

func (h *Handler) MergePR(c *gin.Context) {
	var req schema.MergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "Invalid request body")
		return
	}
	pr, err := h.services.PullRequest.Merge(c.Request.Context(), req.PullRequestID)
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	response := schema.PullRequest{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            schema.PRStatus(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
	c.JSON(http.StatusOK, gin.H{
		"pr": response,
	})
}

func (h *Handler) ReassignReviewer(c *gin.Context) {
	var req schema.ReassignReviewerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "Invalid request body")
		return
	}
	replacedBy, err := h.services.PullRequest.ReassignReviewer(c.Request.Context(), req.PullRequestID, req.OldReviewerID, "")
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	pr, err := h.services.PullRequest.GetByID(c.Request.Context(), req.PullRequestID)
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	response := schema.PullRequest{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            schema.PRStatus(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
	c.JSON(http.StatusOK, gin.H{
		"pr":          response,
		"replaced_by": replacedBy,
	})
}
