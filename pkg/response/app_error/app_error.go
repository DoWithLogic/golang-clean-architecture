package app_error

import "errors"

var (
	ErrInvalidToken              = errors.New("invalid authentication token")
	ErrFailedGetTokenInformation = errors.New("failed to get token information")

	ErrEmailAlreadyExist = errors.New("email already exist")
	ErrInvalidUserType   = errors.New("invalid user_type")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrFailedGenerateJWT = errors.New("failed generate access token")
	ErrInvalidIsActive   = errors.New("invalid is_active")
	ErrStatusValue       = errors.New("status should be 0 or 1")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
