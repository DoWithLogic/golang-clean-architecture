package entities

import (
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
)

type UpdateUsers struct {
	UserID      int64
	Email       string
	Fullname    string
	Password    string
	PhoneNumber string
	UserType    string
	UpdatedAt   time.Time
	UpdatedBy   string
}

func NewUpdateUsers(data dtos.UpdateUserRequest) UpdateUsers {
	return UpdateUsers{
		UserID:      data.UserID,
		Fullname:    data.Fullname,
		PhoneNumber: data.PhoneNumber,
		UserType:    data.UserType,
		UpdatedAt:   time.Now(),
		UpdatedBy:   data.UpdateBy,
	}
}
