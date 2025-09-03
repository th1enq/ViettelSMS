package domain

import "errors"

var (
	ErrUserConflict       = errors.New("username or email already exist")
	ErrInternalServer     = errors.New("interal server error")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)
