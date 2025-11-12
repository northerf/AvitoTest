package domain

import (
	"errors"
	"time"
)

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID                string     `db:"pull_request_id"`
	Name              string     `db:"pull_request_name"`
	AuthorID          string     `db:"author_id"`
	Status            PRStatus   `db:"status"`
	AssignedReviewers []string   `db:"assigned_reviewers"`
	CreatedAt         time.Time  `db:"created_at"`
	MergedAt          *time.Time `db:"merged_at"`
}

func (pr PullRequest) AssignReviewer(reviewerID string) error {
	if pr.Status == StatusMerged {
		return ErrCannotModifyMergedPR
	}
	if pr.AuthorID == reviewerID {
		return ErrCannotReviewOwnPR
	}
	if pr.hasReviewer(reviewerID) {
		return ErrReviewerAlreadyAssigned
	}
	if len(pr.AssignedReviewers) >= 2 {
		return errors.New("cannot assign more than two reviewers")
	}
	pr.AssignedReviewers = append(pr.AssignedReviewers, reviewerID)
	return nil
}

func (pr *PullRequest) RemoveReviewer(reviewerID string) error {
	if pr.Status == StatusMerged {
		return ErrCannotModifyMergedPR
	}
	if !pr.hasReviewer(reviewerID) {
		return ErrReviewerNotAssigned
	}
	newReviewers := make([]string, 0)
	for _, r := range pr.AssignedReviewers {
		if r != reviewerID {
			newReviewers = append(newReviewers, r)
		}
	}
	pr.AssignedReviewers = newReviewers
	return nil
}

func (pr *PullRequest) Merge() error {
	if pr.Status == StatusMerged {
		return nil
	}
	pr.Status = StatusMerged
	now := time.Now()
	pr.MergedAt = &now
	return nil
}

func (pr *PullRequest) hasReviewer(reviewerID string) bool {
	for _, r := range pr.AssignedReviewers {
		if r == reviewerID {
			return true
		}
	}
	return false
}

func (pr *PullRequest) NeedsMoreReviewers() bool {
	return len(pr.AssignedReviewers) < 2
}
