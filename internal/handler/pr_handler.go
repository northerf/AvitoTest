package handler

import (
	"Avito/internal/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) CreatePR(c *gin.Context) {
	var req schema.CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "Invalid request body")
		return
	}
	pr, err := h.services.PullRequest.Create(c.Request.Context(), req.PullRequestName, req.AuthorID)
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	response := schema.PullRequest{ //мейби заменить на функцию
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            schema.PRStatus(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) MergePR(c *gin.Context) {
	prID := c.Query("pull_request_id")
	if prID == "" {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "pull_request_id is required")
		return
	}
	pr, err := h.services.PullRequest.Merge(c.Request.Context(), prID)
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
	c.JSON(http.StatusOK, response)
}

func (h *Handler) ReassignReviewer(c *gin.Context) {
	prID := c.Query("pull_request_id")
	if prID == "" {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "pull_request_id is required")
		return
	}
	var req schema.ReassignReviewerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "Invalid request body")
		return
	}
	err := h.services.PullRequest.ReassignReviewer(c.Request.Context(), prID, req.OldReviewerID, req.NewReviewerID)
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	pr, err := h.services.PullRequest.GetByID(c.Request.Context(), prID)
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
	c.JSON(http.StatusOK, response)
}
