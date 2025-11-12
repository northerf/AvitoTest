package domain

type User struct {
	UserID   string `db:"user_id"`
	Username string `db:"username"`
	TeamName string `db:"team_name"`
	IsActive bool   `db:"is_active"`
}

func (u *User) Deactivate() error {
	if !u.IsActive {
		return ErrUserInactive
	}
	u.IsActive = false
	return nil
}

func (u *User) Activate() error {
	u.IsActive = true
	return nil
}
