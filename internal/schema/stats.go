package schema

type UserStats struct {
	UserID           string `json:"user_id" db:"user_id"`
	Username         string `json:"username" db:"username"`
	ReviewsAssigned  int    `json:"reviews_assigned" db:"reviews_assigned"`
	ReviewsCompleted int    `json:"reviews_completed" db:"reviews_completed"`
}

type PRStats struct {
	PullRequestID   string `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName string `json:"pull_request_name" db:"pull_request_name"`
	AuthorID        string `json:"author_id" db:"author_id"`
	ReviewersCount  int    `json:"reviewers_count" db:"reviewers_count"`
	Status          string `json:"status" db:"status"`
}

type Statistics struct {
	TotalUsers          int         `json:"total_users"`
	TotalActiveUsers    int         `json:"total_active_users"`
	TotalPRs            int         `json:"total_prs"`
	TotalOpenPRs        int         `json:"total_open_prs"`
	TotalMergedPRs      int         `json:"total_merged_prs"`
	TopReviewers        []UserStats `json:"top_reviewers"`
	PRsWithoutReviewers int         `json:"prs_without_reviewers"`
}
