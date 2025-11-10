package db

import "errors"

var (
	ErrUniversityNotFound = errors.New("university not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrScheduleNotFound   = errors.New("schedule not found")
)