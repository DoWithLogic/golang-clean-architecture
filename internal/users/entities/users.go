package entities

import (
	"errors"
	"time"
)

type (
	Users struct {
		UserID      int64
		Fullname    string
		PhoneNumber string
		UserType    string
		IsActive    bool
		CreatedAt   time.Time
		CreatedBy   string
		UpdatedAt   time.Time
		UpdatedBy   string
	}

	LockingOpt struct {
		ForUpdateNoWait bool
		ForUpdate       bool
	}
)

const (
	UserTypeRegular = "regular_user"
	UserTypePremium = "premium_user"
)

var (
	ErrInvalidLockOpt = errors.New("can't do lock with multiple type")
)

func NewCreateUser(data CreateUser) Users {
	return Users{
		Fullname:    data.FullName,
		PhoneNumber: data.PhoneNumber,
		UserType:    UserTypeRegular,
		IsActive:    true,
		CreatedAt:   time.Now(),
		CreatedBy:   "martin",
	}
}
