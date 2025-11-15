package service

import (
	"context"
	"math/rand"
	"time"

	"Avito/internal/domain"
	"Avito/internal/repository"
	"Avito/internal/schema"
)

type TeamServiceImpl struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
	prRepo   repository.PullRequestRepository
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository, prRepo repository.PullRequestRepository) *TeamServiceImpl {
	return &TeamServiceImpl{
		teamRepo: teamRepo,
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (s *TeamServiceImpl) Create(ctx context.Context, teamName string) (*domain.Team, error) {
	team := &domain.Team{
		Name: teamName,
	}
	err := s.teamRepo.Create(ctx, team)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (s *TeamServiceImpl) GetWithMember(ctx context.Context, teamName string) (*domain.Team, []*domain.User, error) {
	return s.teamRepo.GetWithMembers(ctx, teamName)
}

func (s *TeamServiceImpl) AddMember(ctx context.Context, teamName string, userID string) error {
	return s.teamRepo.AddMember(ctx, teamName, userID)
}

func (s *TeamServiceImpl) CreateWithMembers(ctx context.Context, teamName string, members []schema.TeamMember) (*domain.Team, []schema.TeamMember, error) {
	team := &domain.Team{Name: teamName}
	err := s.teamRepo.Create(ctx, team)
	if err != nil {
		return nil, nil, err
	}
	resultMembers := make([]schema.TeamMember, 0, len(members))
	for _, member := range members {
		user := &domain.User{
			UserID:   member.UserID,
			Username: member.Username,
			TeamName: teamName,
			IsActive: member.IsActive,
		}
		err = s.userRepo.Create(ctx, user)
		if err != nil {
			if err == domain.ErrUserExists {
				err = s.userRepo.Update(ctx, user)
				if err != nil {
					return nil, nil, err
				}
			} else {
				return nil, nil, err
			}
		}
		err = s.teamRepo.AddMember(ctx, teamName, member.UserID)
		if err != nil {
			return nil, nil, err
		}
		resultMembers = append(resultMembers, member)
	}
	return team, resultMembers, nil
}

func (s *TeamServiceImpl) DeactivateUsersAndReassign(ctx context.Context, teamName string, userIDs []string) (*schema.TeamDeactivationResult, error) {
	startTime := time.Now()
	err := s.userRepo.DeactivateUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	openPRs, err := s.prRepo.GetOpenPRsByReviewerIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	reassignments := make([]domain.ReviewerReassignment, 0)
	for _, pr := range openPRs {
		candidates, err := s.prRepo.GetActiveReviewersFromTeam(ctx, teamName, pr.AuthorID, 10)
		if err != nil || len(candidates) == 0 {
			continue
		}
		filtered := make([]string, 0)
		for _, c := range candidates {
			isDeactivated := false
			for _, uid := range userIDs {
				if c == uid {
					isDeactivated = true
					break
				}
			}
			if !isDeactivated {
				filtered = append(filtered, c)
			}
		}
		for _, deactivatedID := range userIDs {
			isReviewer := false
			for _, revID := range pr.AssignedReviewers {
				if revID == deactivatedID {
					isReviewer = true
					break
				}
			}
			if isReviewer {
				if len(filtered) > 0 {
					newRev := filtered[rand.Intn(len(filtered))]
					reassignments = append(reassignments, domain.ReviewerReassignment{
						PRID:          pr.ID,
						OldReviewerID: deactivatedID,
						NewReviewerID: newRev,
					})
					filtered = removeFromSlice(filtered, newRev)
				} else {
					reassignments = append(reassignments, domain.ReviewerReassignment{
						PRID:          pr.ID,
						OldReviewerID: deactivatedID,
						NewReviewerID: "",
					})
				}
			}
		}
	}
	err = s.prRepo.BatchReassignReviewers(ctx, reassignments)
	if err != nil {
		return nil, err
	}
	return &schema.TeamDeactivationResult{
		DeactivatedCount: len(userIDs),
		ReassignedPRs:    len(openPRs),
		Duration:         time.Since(startTime),
	}, nil
}

func removeFromSlice(slice []string, item string) []string {
	result := make([]string, 0)
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}
