package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"Avito/internal/domain"
	"Avito/internal/schema"
)

type StatsPostgres struct {
	db *sqlx.DB
}

func NewStatsPostgres(db *sqlx.DB) *StatsPostgres {
	return &StatsPostgres{db: db}
}

func (r *StatsPostgres) GetStatistics(ctx context.Context) (*schema.Statistics, error) {
	stats := &schema.Statistics{}
	err := r.db.GetContext(ctx, &stats.TotalUsers, "SELECT COUNT(*) FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}
	err = r.db.GetContext(ctx, &stats.TotalActiveUsers,
		"SELECT COUNT(*) FROM users WHERE is_active = true")
	if err != nil {
		return nil, fmt.Errorf("failed to get total active users: %w", err)
	}
	err = r.db.GetContext(ctx, &stats.TotalPRs, "SELECT COUNT(*) FROM pull_requests")
	if err != nil {
		return nil, fmt.Errorf("failed to get total PRs: %w", err)
	}
	err = r.db.GetContext(ctx, &stats.TotalOpenPRs,
		"SELECT COUNT(*) FROM pull_requests WHERE status = 'OPEN'")
	if err != nil {
		return nil, fmt.Errorf("failed to get total open PRs: %w", err)
	}
	err = r.db.GetContext(ctx, &stats.TotalMergedPRs,
		"SELECT COUNT(*) FROM pull_requests WHERE status = 'MERGED'")
	if err != nil {
		return nil, fmt.Errorf("failed to get total merged PRs: %w", err)
	}
	err = r.db.GetContext(ctx, &stats.PRsWithoutReviewers, `SELECT COUNT(*) FROM pull_requests pr
        WHERE NOT EXISTS (
            SELECT 1 FROM pr_reviewers prr 
            WHERE prr.pull_request_id = pr.pull_request_id) AND pr.status = 'OPEN'`)
	if err != nil {
		return nil, fmt.Errorf("failed to get PRs without reviewers: %w", err)
	}
	query := `SELECT u.user_id, u.username, COALESCE(COUNT(prr.reviewer_id), 0) as reviews_assigned,
            COALESCE(COUNT(CASE WHEN pr.status = 'MERGED' THEN 1 END), 0) as reviews_completed
        FROM users u
        LEFT JOIN pr_reviewers prr ON u.user_id = prr.reviewer_id
        LEFT JOIN pull_requests pr ON prr.pull_request_id = pr.pull_request_id
        GROUP BY u.user_id, u.username
        HAVING COUNT(prr.reviewer_id) > 0
        ORDER BY reviews_assigned DESC
        LIMIT 10`
	err = r.db.SelectContext(ctx, &stats.TopReviewers, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get top reviewers: %w", err)
	}
	if stats.TopReviewers == nil {
		stats.TopReviewers = []schema.UserStats{}
	}
	return stats, nil
}

func (r *StatsPostgres) GetUserStats(ctx context.Context, userID string) (*schema.UserStats, error) {
	var stats schema.UserStats

	query := `SELECT u.user_id, u.username, COALESCE(COUNT(prr.reviewer_id), 0) as reviews_assigned, COALESCE(COUNT(CASE WHEN pr.status = 'MERGED' THEN 1 END), 0) as reviews_completed
        FROM users u
        LEFT JOIN pr_reviewers prr ON u.user_id = prr.reviewer_id
        LEFT JOIN pull_requests pr ON prr.pull_request_id = pr.pull_request_id
        WHERE u.user_id = $1
        GROUP BY u.user_id, u.username`
	err := r.db.GetContext(ctx, &stats, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	return &stats, nil
}
