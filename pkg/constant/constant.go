package constant

import "errors"

const (
	UserTypeRegular = "regular_user"
	UserTypePremium = "premium_user"

	UserSystem = "SYSTEM"
)

var (
	ErrInvalidLockOpt = errors.New("can't do lock with multiple type")
)

const (
	LengthOfSalt      = 16
	LengthOfRandomKey = 32
)
