package dtos

import (
	"github.com/invopop/validation"
)

type (
	UserLoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	UserLoginResponse struct {
		AccessToken string `json:"access_token"`
		ExpiredAt   int64  `json:"expired_at"`
	}
)

func (ulr UserLoginRequest) Validate() error {
	return validation.ValidateStruct(&ulr,
		validation.Field(&ulr.Email, validation.Required),
		validation.Field(&ulr.Password, validation.Required),
	)
}
