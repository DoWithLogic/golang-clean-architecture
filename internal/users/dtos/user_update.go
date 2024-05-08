package dtos

import (
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
	"github.com/invopop/validation"
)

type UpdateUser struct {
	Fullname    string `json:"fullname"`
	PhoneNumber string `json:"phone_number"`
	UserType    string `json:"user_type"`
}

type UpdateUserRequest struct {
	UserID int64
	UpdateUser
}

func (cup UpdateUser) Validate() error {
	return validation.ValidateStruct(&cup,
		validation.Field(&cup.UserType, validation.In(constant.UserTypePremium, constant.UserTypeRegular)),
	)
}
