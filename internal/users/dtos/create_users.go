package dtos

import (
	"github.com/invopop/validation"
)

type CreateUserPayload struct {
	FullName    string `json:"fullname"`
	PhoneNumber string `json:"phone_number"`
}

type CreateUserResponse struct {
	UserID    int64  `json:"user_id"`
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expired_at"`
}

func (cup CreateUserPayload) Validate() error {
	return validation.ValidateStruct(&cup,
		validation.Field(&cup.FullName, validation.Required, validation.Length(0, 50)),
		validation.Field(&cup.PhoneNumber, validation.Required, validation.Length(0, 13)),
	)
}
