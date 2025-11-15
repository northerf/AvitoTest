package service

import (
	"context"

	"Avito/internal/domain"
	"Avito/internal/repository"
	"Avito/internal/schema"
)

type PullRequestService interface {
	Create(ctx context.Context, id string, name string, authorID string) (*domain.PullRequest, error)
	GetByID(ctx context.Context, id string) (*domain.PullRequest, error)
	Merge(ctx context.Context, id string) (*domain.PullRequest, error)
	AssignReviewer(ctx context.Context, prID string, count int) error
	ReassignReviewer(ctx context.Context, prID string, oldReviewerID string, newReviewerID string) (string, error)
	GetByReviewerID(ctx context.Context, reviewerID string) ([]*domain.PullRequest, error)
}

type UserService interface {
	Create(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	SetActive(ctx context.Context, userID string, isActive bool) error
}

type TeamService interface {
	Create(ctx context.Context, teamName string) (*domain.Team, error)
	GetWithMember(ctx context.Context, teamName string) (*domain.Team, []*domain.User, error)
	AddMember(ctx context.Context, teamName string, userID string) error
	CreateWithMembers(ctx context.Context, teamName string, members []schema.TeamMember) (*domain.Team, []schema.TeamMember, error)
	DeactivateUsersAndReassign(ctx context.Context, teamName string, userIDs []string) (*schema.TeamDeactivationResult, error)
}

type StatsService interface {
	GetStatistics(ctx context.Context) (*schema.Statistics, error)
	GetUserStats(ctx context.Context, userID string) (*schema.UserStats, error)
}

type Service struct {
	PullRequest PullRequestService
	User        UserService
	Team        TeamService
	Stats       StatsService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		PullRequest: NewPullRequestService(repo.PullRequest, repo.User, repo.Team),
		User:        NewUserService(repo.User),
		Team:        NewTeamService(repo.Team, repo.User, repo.PullRequest),
		Stats:       NewStatsService(repo.Stats),
	}
}
