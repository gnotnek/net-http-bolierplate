package apperror

import "errors"

var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrInvalidPassword  = errors.New("invalid password")
)
