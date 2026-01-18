package domain

import "errors"

var (
	ErrInternal = errors.New("internal server error")
	ErrNotFound = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
)
