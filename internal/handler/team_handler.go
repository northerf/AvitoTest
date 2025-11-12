package handler

import (
	"Avito/internal/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) CreateTeam(c *gin.Context) {
	var req schema.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "Invalid request body")
		return
	}
	team, err := h.services.Team.Create(c.Request.Context(), req.TeamName)
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, schema.Team{
		TeamName: team.Name,
		Members:  []schema.TeamMember{},
	})
}

func (h *Handler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		ErrorResponse(c, http.StatusBadRequest, schema.ErrNotFound, "team_name is required")
		return
	}
	team, users, err := h.services.Team.GetWithMember(c.Request.Context(), teamName)
	if err != nil {
		MapErrorToResponse(c, err)
		return
	}
	members := make([]schema.TeamMember, len(users))
	for i, user := range users {
		members[i] = schema.TeamMember{
			UserID:   user.UserID,
			Username: user.Username,
			IsActive: user.IsActive,
		}
	}
	c.JSON(http.StatusOK, schema.Team{
		TeamName: team.Name,
		Members:  members,
	})
}
