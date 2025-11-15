package schema

import "time"

type DeactivateUsersRequest struct {
	TeamName string   `json:"team_name" binding:"required"`
	UserIDs  []string `json:"user_ids" binding:"required,min=1"`
}

type TeamDeactivationResult struct {
	DeactivatedCount int           `json:"deactivated_count"`
	ReassignedPRs    int           `json:"reassigned_prs"`
	Duration         time.Duration `json:"-"`
}

type DeactivateUsersResponse struct {
	DeactivatedCount int   `json:"deactivated_count"`
	ReassignedPRs    int   `json:"reassigned_prs"`
	DurationMs       int64 `json:"duration_ms"`
}
