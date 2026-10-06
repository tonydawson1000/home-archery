package domain

import "errors"

var (
	ErrInvalidScoreCode = errors.New("invalid_score_code")
	ErrSessionCompleted = errors.New("session_completed")
	ErrEmptyComplete    = errors.New("empty_complete")
)
