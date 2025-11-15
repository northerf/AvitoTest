package schema

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type CreateTeamRequest struct {
	TeamName string       `json:"team_name" binding:"required"`
	Members  []TeamMember `json:"members" binding:"required"`
}
