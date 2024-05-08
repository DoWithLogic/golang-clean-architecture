package entities

import (
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_crypto"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
)

type (
	User struct {
		UserID      int64
		Email       string
		Password    string
		Fullname    string
		PhoneNumber string
		UserType    string
		IsActive    bool
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	LockingOpt struct {
		PessimisticLocking bool
	}
)

func NewCreateUser(data dtos.CreateUserRequest, cfg config.Config) User {
	return User{
		Fullname:    data.FullName,
		Email:       data.Email,
		Password:    app_crypto.NewCrypto(cfg.Authentication.Key).EncodeSHA256(data.Password),
		PhoneNumber: data.PhoneNumber,
		UserType:    constant.UserTypeRegular,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}
}
