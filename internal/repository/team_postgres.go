package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"Avito/internal/domain"
)

type TeamPostgres struct {
	db *sqlx.DB
}

func NewTeamPostgres(db *sqlx.DB) *TeamPostgres {
	return &TeamPostgres{db: db}
}

func (r *TeamPostgres) Create(ctx context.Context, team *domain.Team) error {
	query := `INSERT INTO teams (team_name) VALUES ($1)`
	_, err := r.db.ExecContext(ctx, query, team.Name)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return domain.ErrTeamExists
		}
		return fmt.Errorf("failed to create team: %w", err)
	}
	return nil
}

func (r *TeamPostgres) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	query := `SELECT team_name FROM teams WHERE team_name = $1`
	var team domain.Team
	err := r.db.GetContext(ctx, &team, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTeamNotFound
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	return &team, nil
}

func (r *TeamPostgres) GetWithMembers(ctx context.Context, teamName string) (*domain.Team, []*domain.User, error) {
	team, err := r.GetByName(ctx, teamName)
	if err != nil {
		return nil, nil, err
	}
	query := `SELECT u.user_id, u.username, u.team_name, u.is_active FROM users u
        INNER JOIN team_members tm ON u.user_id = tm.user_id
        WHERE tm.team_name = $1
        ORDER BY u.username`
	var users []*domain.User
	err = r.db.SelectContext(ctx, &users, query, teamName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get team members: %w", err)
	}
	return team, users, nil
}

func (r *TeamPostgres) AddMember(ctx context.Context, teamName string, userID string) error {
	query := `INSERT INTO team_members (team_name, user_id) VALUES ($1, $2) ON CONFLICT (team_name, user_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, teamName, userID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return domain.ErrTeamNotFound
		}
		return fmt.Errorf("failed to add team member: %w", err)
	}
	updateQuery := `UPDATE users SET team_name = $1 WHERE user_id = $2`
	_, err = r.db.ExecContext(ctx, updateQuery, teamName, userID)
	if err != nil {
		return fmt.Errorf("failed to update user team: %w", err)
	}
	return nil
}

func (r *TeamPostgres) GetByUserID(ctx context.Context, userID string) (*domain.Team, error) {
	query := `SELECT t.team_name FROM teams t
        INNER JOIN team_members tm ON t.team_name = tm.team_name
        WHERE tm.user_id = $1`
	var team domain.Team
	err := r.db.GetContext(ctx, &team, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTeamNotFound
		}
		return nil, fmt.Errorf("failed to get team by user: %w", err)
	}
	return &team, nil
}
