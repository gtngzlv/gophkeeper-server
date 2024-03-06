package errors

import "errors"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrFailedGetUserID    = errors.New("failed to get userID from context")
	ErrFailedSaveData     = errors.New("failed to saved data")
	ErrFailedInsertData   = errors.New("failed to insert data")
)
