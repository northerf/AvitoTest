package service

import (
	"Avito/internal/domain"
	"Avito/internal/repository"
	"context"
	"fmt"
)

type TeamServiceImpl struct {
	teamRepo repository.TeamRepository
}

func NewTeamService(teamRepo repository.TeamRepository) *TeamServiceImpl {
	return &TeamServiceImpl{
		teamRepo: teamRepo,
	}
}

func (s *TeamServiceImpl) Create(ctx context.Context, teamName string) (*domain.Team, error) {
	team := &domain.Team{
		Name: teamName,
	}
	err := s.teamRepo.Create(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}
	return team, nil
}

func (s *TeamServiceImpl) GetWithMember(ctx context.Context, teamName string) (*domain.Team, []*domain.User, error) {
	return s.teamRepo.GetWithMembers(ctx, teamName)
}

func (s *TeamServiceImpl) AddMember(ctx context.Context, teamName string, userID string) error {
	return s.teamRepo.AddMember(ctx, teamName, userID)
}
