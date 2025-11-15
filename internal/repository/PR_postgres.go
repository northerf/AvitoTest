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

type PullRequestPostgres struct {
	db *sqlx.DB
}

func NewPullRequestPostgres(db *sqlx.DB) *PullRequestPostgres {
	return &PullRequestPostgres{db: db}
}

func (r *PullRequestPostgres) Create(ctx context.Context, pr *domain.PullRequest) error {
	query := `INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at) VALUES ($1, $2, $3, $4, NOW())
        RETURNING created_at`
	err := r.db.QueryRowContext(ctx, query, pr.ID, pr.Name, pr.AuthorID, pr.Status).Scan(&pr.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.ErrPRExists
		}
		return fmt.Errorf("failed to create pull request: %w", err)
	}
	return nil
}

func (r *PullRequestPostgres) GetByID(ctx context.Context, id string) (*domain.PullRequest, error) {
	query := `SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at,
            COALESCE(array_agg(prr.reviewer_id) FILTER (WHERE prr.reviewer_id IS NOT NULL), '{}') as assigned_reviewers
        FROM pull_requests pr
        LEFT JOIN pr_reviewers prr ON pr.pull_request_id = prr.pull_request_id
        WHERE pr.pull_request_id = $1
        GROUP BY pr.pull_request_id`
	var pr domain.PullRequest
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
		pq.Array(&pr.AssignedReviewers),
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrPRNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}
	return &pr, nil
}

func (r *PullRequestPostgres) Merge(ctx context.Context, id string) (*domain.PullRequest, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	var pr domain.PullRequest
	query := `SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
        FROM pull_requests
		WHERE pull_request_id = $1
        FOR UPDATE`
	err = tx.GetContext(ctx, &pr, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPRNotFound
		}
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}
	if pr.Status == domain.StatusMerged {
		return r.GetByID(ctx, id)
	}
	updateQuery := `UPDATE pull_requests SET status = $1, merged_at = NOW()
        WHERE pull_request_id = $2
        RETURNING merged_at`
	err = tx.QueryRowContext(ctx, updateQuery, domain.StatusMerged, id).Scan(&pr.MergedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to merge pull request: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return r.GetByID(ctx, id)
}

func (r *PullRequestPostgres) AssignReviewer(ctx context.Context, prID string, reviewerID string) error {
	var status string
	checkQuery := `SELECT status FROM pull_requests WHERE pull_request_id = $1`
	err := r.db.GetContext(ctx, &status, checkQuery, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrPRNotFound
		}
		return fmt.Errorf("failed to check PR status: %w", err)
	}
	if status == string(domain.StatusMerged) {
		return domain.ErrCannotModifyMergedPR
	}
	var isActive bool
	activeQuery := `SELECT is_active FROM users WHERE user_id = $1`
	err = r.db.GetContext(ctx, &isActive, activeQuery, reviewerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("failed to check reviewer status: %w", err)
	}
	if !isActive {
		return domain.ErrUserInactive
	}
	query := `INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)
        ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING`
	_, err = r.db.ExecContext(ctx, query, prID, reviewerID)
	if err != nil {
		return fmt.Errorf("failed to assign reviewer: %w", err)
	}
	return nil
}

func (r *PullRequestPostgres) RemoveReviewer(ctx context.Context, prID string, reviewerID string) error {
	var status string
	checkQuery := `SELECT status FROM pull_requests WHERE pull_request_id = $1`
	err := r.db.GetContext(ctx, &status, checkQuery, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrPRNotFound
		}
		return fmt.Errorf("failed to check PR status: %w", err)
	}
	if status == string(domain.StatusMerged) {
		return domain.ErrCannotModifyMergedPR
	}
	query := `DELETE FROM pr_reviewers WHERE pull_request_id = $1 AND reviewer_id = $2`
	result, err := r.db.ExecContext(ctx, query, prID, reviewerID)
	if err != nil {
		return fmt.Errorf("failed to remove reviewer: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrReviewerNotAssigned
	}
	return nil
}

func (r *PullRequestPostgres) GetByReviewerID(ctx context.Context, reviewerID string) ([]*domain.PullRequest, error) {
	query := `SELECT DISTINCT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at
        FROM pull_requests pr
        INNER JOIN pr_reviewers prr ON pr.pull_request_id = prr.pull_request_id
        WHERE prr.reviewer_id = $1
        ORDER BY pr.created_at DESC`
	var prs []*domain.PullRequest
	err := r.db.SelectContext(ctx, &prs, query, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull requests by reviewer: %w", err)
	}
	for _, pr := range prs {
		reviewersQuery := `SELECT reviewer_id FROM pr_reviewers WHERE pull_request_id = $1`
		err = r.db.SelectContext(ctx, &pr.AssignedReviewers, reviewersQuery, pr.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get reviewers: %w", err)
		}
	}
	return prs, nil
}

func (r *PullRequestPostgres) GetActiveReviewersFromTeam(ctx context.Context, teamName string, excludeUserID string, limit int) ([]string, error) {
	query := `SELECT u.user_id FROM users u
        INNER JOIN team_members tm ON u.user_id = tm.user_id
        WHERE tm.team_name = $1
          AND u.is_active = true AND u.user_id != $2
        ORDER BY RANDOM()
        LIMIT $3`
	var reviewers []string
	err := r.db.SelectContext(ctx, &reviewers, query, teamName, excludeUserID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get active reviewers: %w", err)
	}
	return reviewers, nil
}

func (r *PullRequestPostgres) GetActiveReviewersFromUserTeam(ctx context.Context, userID string, excludeUserIDs []string, limit int) ([]string, error) {
	query := `SELECT u.user_id FROM users u
        INNER JOIN team_members tm ON u.user_id = tm.user_id
        WHERE tm.team_name = (SELECT team_name FROM users WHERE user_id = $1) 
        AND u.user_id != ALL($2) AND u.is_active = true
        ORDER BY RANDOM()
        LIMIT $3`
	var reviewers []string
	err := r.db.SelectContext(ctx, &reviewers, query, userID, pq.Array(excludeUserIDs), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get active reviewers from user team: %w", err)
	}
	return reviewers, nil
}

func (r *PullRequestPostgres) GetOpenPRsByReviewerIDs(ctx context.Context, reviewerIDs []string) ([]*domain.PullRequest, error) {
	if len(reviewerIDs) == 0 {
		return []*domain.PullRequest{}, nil
	}
	query := `SELECT DISTINCT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at,pr.merged_at
        FROM pull_requests pr
        INNER JOIN pr_reviewers prr ON pr.pull_request_id = prr.pull_request_id
        WHERE prr.reviewer_id = ANY($1) AND pr.status = 'OPEN'`
	var prs []*domain.PullRequest
	err := r.db.SelectContext(ctx, &prs, query, pq.Array(reviewerIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to get open PRs: %w", err)
	}
	for _, pr := range prs {
		reviewersQuery := `SELECT reviewer_id FROM pr_reviewers WHERE pull_request_id = $1`
		err = r.db.SelectContext(ctx, &pr.AssignedReviewers, reviewersQuery, pr.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get reviewers for PR %s: %w", pr.ID, err)
		}
	}
	return prs, nil
}

func (r *PullRequestPostgres) BatchReassignReviewers(ctx context.Context, reassignments []domain.ReviewerReassignment) error {
	if len(reassignments) == 0 {
		return nil
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	for _, reassignment := range reassignments {
		deleteQuery := `DELETE FROM pr_reviewers WHERE pull_request_id = $1 AND reviewer_id = $2`
		_, err := tx.ExecContext(ctx, deleteQuery, reassignment.PRID, reassignment.OldReviewerID)
		if err != nil {
			return fmt.Errorf("failed to remove old reviewer: %w", err)
		}
		if reassignment.NewReviewerID != "" {
			insertQuery := `INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2) ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING`
			_, err = tx.ExecContext(ctx, insertQuery, reassignment.PRID, reassignment.NewReviewerID)
			if err != nil {
				return fmt.Errorf("failed to assign new reviewer: %w", err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch reassignment: %w", err)
	}
	return nil
}
