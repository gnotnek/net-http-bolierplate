package apperror

import "errors"

var (
	ErrUserDeleted         = errors.New("user already deleted")
	ErrUserInvalidPassword = errors.New("invalid password")
	ErrNotFound            = errors.New("entity not found")
)
