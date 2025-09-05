package dtos

import (
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
	"github.com/invopop/validation"
)

type (
	UserLoginRequest struct {
		ContactType  types.CONTACT_TYPE `json:"contact_type"`
		ContactValue string             `json:"contact_value"`
		Password     string             `json:"password"`
	}

	UserLoginResponse struct {
		AccessToken string `json:"access_token"`
		ExpiredAt   int64  `json:"expired_at"`
	}
)

func (ulr UserLoginRequest) Validate() error {
	return validation.ValidateStruct(&ulr,
		validation.Field(&ulr.ContactType, validation.Required),
		validation.Field(&ulr.ContactValue, validation.Required),
		validation.Field(&ulr.Password, validation.Required),
	)
}

func ToUserLoginResponse(accessToken string, expiredAt int64) UserLoginResponse {
	return UserLoginResponse{
		AccessToken: accessToken,
		ExpiredAt:   expiredAt,
	}
}
