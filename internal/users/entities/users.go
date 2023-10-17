package entities

import (
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/utils"
)

type (
	Users struct {
		UserID      int64
		Email       string
		Password    string
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
		PessimisticLocking bool
	}
)

func NewCreateUser(data dtos.CreateUserRequest, cfg config.Config) Users {
	return Users{
		Fullname:    data.FullName,
		Email:       data.Email,
		Password:    utils.Encrypt(data.Password, cfg),
		PhoneNumber: data.PhoneNumber,
		UserType:    constant.UserTypeRegular,
		IsActive:    true,
		CreatedAt:   time.Now(),
		CreatedBy:   "SYSTEM",
	}
}
