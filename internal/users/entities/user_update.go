package entities

import (
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
)

type UpdateUser struct {
	UserID      int64
	Fullname    string
	PhoneNumber string
	UserType    string
	UpdatedAt   time.Time
}

func NewUpdateUser(data dtos.UpdateUserRequest) UpdateUser {
	return UpdateUser{
		UserID:      data.UserID,
		Fullname:    data.Fullname,
		PhoneNumber: data.PhoneNumber,
		UserType:    data.UserType,
		UpdatedAt:   time.Now(),
	}
}
