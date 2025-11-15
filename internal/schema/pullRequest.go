package schema

import "time"

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            PRStatus   `json:"status" binding:"oneof=OPEN MERGED"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

type CreatePRRequest struct {
	PullRequestName string `json:"pull_request_name" binding:"required"`
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

type AssignReviewersRequest struct {
	Count int `json:"count" binding:"required,gte=0,lte=2"`
}

type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldReviewerID string `json:"old_reviewer_id" binding:"required"`
	NewReviewerID string `json:"new_reviewer_id"`
}
