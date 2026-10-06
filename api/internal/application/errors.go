package application

import "errors"

var (
	ErrArcherNotFound       = errors.New("archer_not_found")
	ErrSessionNotFound      = errors.New("session_not_found")
	ErrPersonalBestNotFound = errors.New("personal_best_not_found")
)
