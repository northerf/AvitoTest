package handler

import (
	"github.com/gin-gonic/gin"

	"Avito/internal/domain"
	"Avito/internal/schema"
)

func ErrorResponse(c *gin.Context, statusCode int, code schema.ErrorCode, message string) {
	c.JSON(statusCode, schema.ErrorResponse{
		Error: schema.ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func MapErrorToResponse(c *gin.Context, err error) {
	switch err {
	case domain.ErrTeamExists:
		ErrorResponse(c, 400, schema.ErrTeamExists, "Team already exists")
	case domain.ErrTeamNotFound:
		ErrorResponse(c, 404, schema.ErrNotFound, "Team not found")
	case domain.ErrUserNotFound:
		ErrorResponse(c, 404, schema.ErrNotFound, "User not found")
	case domain.ErrUserExists:
		ErrorResponse(c, 400, schema.ErrNotFound, "User already exists")
	case domain.ErrPRExists:
		ErrorResponse(c, 409, schema.ErrPRExists, "PR already exists")
	case domain.ErrPRNotFound:
		ErrorResponse(c, 404, schema.ErrNotFound, "PR not found")
	case domain.ErrCannotModifyMergedPR:
		ErrorResponse(c, 409, schema.ErrPRMerged, "Cannot modify merged PR")
	case domain.ErrReviewerNotAssigned:
		ErrorResponse(c, 409, schema.ErrNotAssigned, "Reviewer not assigned")
	case domain.ErrNoActiveCandidates:
		ErrorResponse(c, 409, schema.ErrNoCandidate, "No active candidates available")
	default:
		ErrorResponse(c, 500, schema.ErrNotFound, "Internal server error")
	}
}
