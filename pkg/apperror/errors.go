package apperror

import "errors"

var (
	ErrEmailAlreadyExist = errors.New("email already exist")
	ErrInvalidUserType   = errors.New("invalid user_type")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrFailedGenerateJWT = errors.New("failed generate access token")
	ErrInvalidIsActive   = errors.New("invalid is_active")
	ErrStatusValue       = errors.New("status should be 0 or 1")
)
