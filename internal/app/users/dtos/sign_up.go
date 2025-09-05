package dtos

import (
	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
	"github.com/invopop/validation"
)

type SignUpRequest struct {
	Name         string             `json:"name"`
	ContactType  types.CONTACT_TYPE `json:"contact_type"`
	ContactValue string             `json:"contact_value"`
	Password     string             `json:"password"`
}

func (s SignUpRequest) ToUserEntity(encryptedPassword string) *entities.User {
	return &entities.User{
		Name:         s.Name,
		ContactType:  s.ContactType,
		ContactValue: s.ContactValue,
		Password:     encryptedPassword,
		Status:       types.PENDING,
	}
}

func (v SignUpRequest) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Name, validation.Required),
		validation.Field(&v.Password, validation.Required),
		validation.Field(&v.ContactType, validation.Required, validation.In(types.CONTACT_TYPE_EMAIL, types.CONTACT_TYPE_PHONE)),
		validation.Field(&v.ContactValue, validation.Required),
	)
}
