package repository

import (
	"Avito/internal/domain"
	"context"
	"github.com/jmoiron/sqlx"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pr *domain.PullRequest) error
	GetByID(ctx context.Context, id string) (*domain.PullRequest, error)
	Merge(ctx context.Context, id string) (*domain.PullRequest, error)
	GetByReviewerID(ctx context.Context, reviewerID string) ([]*domain.PullRequest, error)
	AssignReviewer(ctx context.Context, prID string, reviewerID string) error
	RemoveReviewer(ctx context.Context, prID string, reviewerID string) error
	GetActiveReviewersFromTeam(ctx context.Context, teamName string, excludeUserID string, limit int) ([]string, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	SetActive(ctx context.Context, userID string, isActive bool) error
	GetByTeamName(ctx context.Context, teamName string) ([]*domain.User, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByName(ctx context.Context, name string) (*domain.Team, error)
	GetWithMembers(ctx context.Context, teamName string) (*domain.Team, []*domain.User, error)
	AddMember(ctx context.Context, teamName string, userID string) error
	GetByUserID(ctx context.Context, userID string) (*domain.Team, error)
}

type Repository struct {
	PullRequest PullRequestRepository
	User        UserRepository
	Team        TeamRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		PullRequest: NewPullRequestPostgres(db),
		User:        NewUserPostgres(db),
		Team:        NewTeamPostgres(db),
	}
}
