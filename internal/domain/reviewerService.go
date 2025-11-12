package domain

import "math/rand"

type ReviewerService struct{}

func NewReviewerService() *ReviewerService {
	return &ReviewerService{}
}

func (s *ReviewerService) SelectReviewers(team *Team, authorID string, count int) ([]string, error) {
	activeMembers := team.GetActiveMembers()
	candidates := make([]User, 0)
	for _, member := range activeMembers {
		if member.UserID != authorID {
			candidates = append(candidates, member)
		}
	}
	if len(candidates) < count {
		count = len(candidates)
	}
	selected := make([]string, 0, count)
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for i := 0; i < count && i < len(candidates); i++ {
		selected = append(selected, candidates[i].UserID)
	}
	return selected, nil
}
