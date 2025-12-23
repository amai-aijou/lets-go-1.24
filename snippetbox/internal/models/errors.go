package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// 10.2: Use this error if a user tries to log in with incorrect email/password
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// 10.2: Use this error if a user tries signing up with an email already in use
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
