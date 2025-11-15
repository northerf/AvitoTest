package domain

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrUserExists              = errors.New("user already exists")
	ErrUserInactive            = errors.New("user is inactive")
	ErrTeamNotFound            = errors.New("team not found")
	ErrTeamExists              = errors.New("team already exists")
	ErrPRNotFound              = errors.New("pull request not found")
	ErrPRExists                = errors.New("pull request already exists")
	ErrCannotModifyMergedPR    = errors.New("cannot modify merged pull request")
	ErrReviewerNotAssigned     = errors.New("reviewer not assigned")
	ErrReviewerAlreadyAssigned = errors.New("reviewer already assigned")
	ErrNoActiveCandidates      = errors.New("no active candidates available")
	ErrCannotReviewOwnPR       = errors.New("cannot review own pull request")
)
