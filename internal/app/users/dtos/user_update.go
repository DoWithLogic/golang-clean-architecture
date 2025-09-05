package dtos

import (
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
)

type UserUpdateRequest struct {
	ID int64 `param:"id"`
	UserUpdate
}

type UserUpdate struct {
	Name         *string             `json:"name"`
	ContactType  *types.CONTACT_TYPE `json:"contact_type"`
	ContactValue *string             `json:"contact_value"`
	BirthDate    *string             `json:"birth_date"`
	Language     *types.LANGUAGE     `json:"language"`
	Password     *string             `json:"password"`
}

func (u UserUpdateRequest) ToUpdateUserEntity(encryptedPassword *string) *entities.UpdateUser {
	return &entities.UpdateUser{
		ID:           u.ID,
		Name:         u.Name,
		ContactType:  u.ContactType,
		ContactValue: u.ContactValue,
		BirthDate:    u.BirthDate,
		Language:     u.Language,
		Password:     encryptedPassword,
		UpdatedAt:    time.Now(),
	}
}
