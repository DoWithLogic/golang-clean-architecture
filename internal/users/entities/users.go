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
		CreatedAt   string
		CreatedBy   string
		UpdatedAt   string
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

func (locking *LockingOpt) Validate() error {
	if locking.ForUpdate && locking.ForUpdateNoWait {
		return ErrInvalidLockOpt
	}

	return nil
}

func NewUser(data Users) *Users {
	return &Users{
		Fullname:    data.Fullname,
		PhoneNumber: data.PhoneNumber,
		UserType:    UserTypeRegular,
		IsActive:    true,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		CreatedBy:   "martin",
	}
}
