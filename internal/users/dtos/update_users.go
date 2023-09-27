package dtos

import (
	"errors"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	"github.com/invopop/validation"
)

type UpdateUserPayload struct {
	Fullname    string `json:"fullname"`
	PhoneNumber string `json:"phone_number"`
	UserType    string `json:"user_type"`
}

var (
	ErrInvalidUserType = errors.New("invalid user_type")
)

func (cup UpdateUserPayload) Validate() error {
	var validationFields []*validation.FieldRules

	if cup.UserType != "" && (cup.UserType != entities.UserTypePremium && cup.UserType != entities.UserTypeRegular) {
		return ErrInvalidUserType
	}

	if cup.Fullname != "" {
		validationFields = append(validationFields, validation.Field(&cup.Fullname, validation.Required, validation.Length(0, 50)))
	}

	if cup.PhoneNumber != "" {
		validationFields = append(validationFields, validation.Field(&cup.PhoneNumber, validation.Required, validation.Length(0, 13)))
	}

	return validation.ValidateStruct(&cup, validationFields...)
}
