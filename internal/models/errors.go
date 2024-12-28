package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// If user tries to login with an incorrect email address/password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// If user tries to signup with an already registerred email.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
