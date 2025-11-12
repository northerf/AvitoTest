package schema

type User struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
}

type UpdateUserActivityRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}
