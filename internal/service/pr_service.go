package service

import (
	"Avito/internal/domain"
	"Avito/internal/repository"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type PullRequestServiceImpl struct {
	prRepo          repository.PullRequestRepository
	teamRepo        repository.TeamRepository
	userRepo        repository.UserRepository
	reviewerService *domain.ReviewerService
}

func NewPullRequestService(prRepo repository.PullRequestRepository,
	userRepo repository.UserRepository, teamRepo repository.TeamRepository) *PullRequestServiceImpl {
	return &PullRequestServiceImpl{
		prRepo:          prRepo,
		userRepo:        userRepo,
		teamRepo:        teamRepo,
		reviewerService: domain.NewReviewerService(),
	}
}

func (s *PullRequestServiceImpl) Create(ctx context.Context, name string, authorID string) (*domain.PullRequest, error) {
	_, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get author: %w", err)
	}
	_, err = s.teamRepo.GetByUserID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get author team: %w", err)
	}
	pr := &domain.PullRequest{
		ID:                uuid.New().String(),
		Name:              name,
		AuthorID:          authorID,
		Status:            domain.StatusOpen,
		AssignedReviewers: []string{},
		CreatedAt:         time.Now(),
	}
	err = s.prRepo.Create(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to create PR: %w", err)
	}
	err = s.AssignReviewer(ctx, pr.ID, 2)
	if err != nil {
		fmt.Printf("failed to assign reviewers: %v\n", err)
	}
	return s.prRepo.GetByID(ctx, pr.ID)
}

func (s *PullRequestServiceImpl) AssignReviewer(ctx context.Context, prID string, count int) error {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return fmt.Errorf("failed to get PR: %w", err)
	}
	team, err := s.teamRepo.GetByUserID(ctx, pr.AuthorID)
	if err != nil {
		return fmt.Errorf("failed to get team: %w", err)
	}
	reviewers, err := s.prRepo.GetActiveReviewersFromTeam(ctx, team.Name, pr.AuthorID, count)
	if err != nil {
		return fmt.Errorf("failed to get active reviewers: %w", err)
	}
	for _, reviewerID := range reviewers {
		err = pr.AssignReviewer(reviewerID)
		if err != nil {
			continue
		}
		err = s.prRepo.AssignReviewer(ctx, prID, reviewerID)
		if err != nil {
			return fmt.Errorf("failed to assign reviewer: %w", err)
		}
	}
	return nil
}

func (s *PullRequestServiceImpl) ReassignReviewer(ctx context.Context, prID string, oldReviewerID string, newReviewerID string) error {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return fmt.Errorf("failed to get PR: %w", err)
	}
	err = pr.RemoveReviewer(oldReviewerID)
	if err != nil {
		return err
	}
	err = pr.AssignReviewer(newReviewerID)
	if err != nil {
		return err
	}
	err = s.prRepo.RemoveReviewer(ctx, prID, oldReviewerID)
	if err != nil {
		return fmt.Errorf("failed to remove reviewer: %w", err)
	}
	err = s.prRepo.AssignReviewer(ctx, prID, newReviewerID)
	if err != nil {
		return fmt.Errorf("failed to assign reviewer: %w", err)
	}
	return nil
}

func (s *PullRequestServiceImpl) Merge(ctx context.Context, id string) (*domain.PullRequest, error) {
	return s.prRepo.Merge(ctx, id)
}

func (s *PullRequestServiceImpl) GetByID(ctx context.Context, id string) (*domain.PullRequest, error) {
	return s.prRepo.GetByID(ctx, id)
}

func (s *PullRequestServiceImpl) GetByReviewerID(ctx context.Context, reviewerID string) ([]*domain.PullRequest, error) {
	return s.prRepo.GetByReviewerID(ctx, reviewerID)
}
