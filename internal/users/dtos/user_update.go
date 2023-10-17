package dtos

import (
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
)

type UpdateUserRequest struct {
	UserID      int64  `json:"-"`
	Fullname    string `json:"fullname"`
	PhoneNumber string `json:"phone_number"`
	UserType    string `json:"user_type"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

var ()

func (cup UpdateUserRequest) Validate() error {
	if cup.UserType != "" && cup.UserType != constant.UserTypePremium && cup.UserType != constant.UserTypeRegular {
		return apperror.ErrInvalidUserType
	}

	return nil
}
