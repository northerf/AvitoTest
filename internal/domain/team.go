package domain

type Team struct {
	Name    string `db:"team_name"`
	Members []User
}

func (t *Team) GetActiveMembers() []User {
	active := make([]User, 0)
	for _, member := range t.Members {
		if member.IsActive {
			active = append(active, member)
		}
	}
	return active
}

func (t *Team) HasMember(userID string) bool {
	for _, member := range t.Members {
		if member.UserID == userID {
			return true
		}
	}
	return false
}
