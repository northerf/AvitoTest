package service

import (
	"context"
	"fmt"
	"time"

	"Avito/internal/domain"
	"Avito/internal/repository"
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

func (s *PullRequestServiceImpl) Create(ctx context.Context, id string, name string, authorID string) (*domain.PullRequest, error) {
	_, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, err
	}
	_, err = s.teamRepo.GetByUserID(ctx, authorID)
	if err != nil {
		return nil, err
	}
	pr := &domain.PullRequest{
		ID:                id,
		Name:              name,
		AuthorID:          authorID,
		Status:            domain.StatusOpen,
		AssignedReviewers: []string{},
		CreatedAt:         time.Now(),
	}
	err = s.prRepo.Create(ctx, pr)
	if err != nil {
		return nil, err
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
		return err
	}
	team, err := s.teamRepo.GetByUserID(ctx, pr.AuthorID)
	if err != nil {
		return err
	}
	reviewers, err := s.prRepo.GetActiveReviewersFromTeam(ctx, team.Name, pr.AuthorID, count)
	if err != nil {
		return err
	}
	for _, reviewerID := range reviewers {
		err = pr.AssignReviewer(reviewerID)
		if err != nil {
			continue
		}
		err = s.prRepo.AssignReviewer(ctx, prID, reviewerID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PullRequestServiceImpl) ReassignReviewer(ctx context.Context, prID string, oldReviewerID string, newReviewerID string) (string, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return "", err
	}
	err = pr.RemoveReviewer(oldReviewerID)
	if err != nil {
		return "", err
	}
	var selectedReviewerID string
	if newReviewerID == "" {
		excludeIDs := []string{pr.AuthorID, oldReviewerID}
		for _, rid := range pr.AssignedReviewers {
			if rid != oldReviewerID {
				excludeIDs = append(excludeIDs, rid)
			}
		}
		reviewers, err := s.prRepo.GetActiveReviewersFromUserTeam(ctx, oldReviewerID, excludeIDs, 1)
		if err != nil {
			return "", err
		}
		if len(reviewers) == 0 {
			return "", domain.ErrNoActiveCandidates
		}
		selectedReviewerID = reviewers[0]
	} else {
		selectedReviewerID = newReviewerID
	}
	err = pr.AssignReviewer(selectedReviewerID)
	if err != nil {
		return "", err
	}
	err = s.prRepo.RemoveReviewer(ctx, prID, oldReviewerID)
	if err != nil {
		return "", err
	}
	err = s.prRepo.AssignReviewer(ctx, prID, selectedReviewerID)
	if err != nil {
		return "", err
	}
	return selectedReviewerID, nil
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
