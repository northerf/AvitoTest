package schema

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type AddTeamMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
}
