package repository

import (
	"Avito/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (user_id, username, team_name, is_active)
        VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, user.UserID, user.Username, user.TeamName, user.IsActive)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return domain.ErrUserExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserPostgres) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1`
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *UserPostgres) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT user_id, username, team_name, is_active FROM users WHERE username = $1`
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *UserPostgres) SetActive(ctx context.Context, userID string, isActive bool) error {
	query := `UPDATE users SET is_active = $1 WHERE user_id = $2`
	result, err := r.db.ExecContext(ctx, query, isActive, userID)
	if err != nil {
		return fmt.Errorf("failed to set user active status: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *UserPostgres) GetByTeamName(ctx context.Context, teamName string) ([]*domain.User, error) {
	query := `SELECT u.user_id, u.username, u.team_name, u.is_active
        FROM users u
        INNER JOIN team_members tm ON u.user_id = tm.user_id
        WHERE tm.team_name = $1
        ORDER BY u.username`
	var users []*domain.User
	err := r.db.SelectContext(ctx, &users, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by team: %w", err)
	}
	return users, nil
}
