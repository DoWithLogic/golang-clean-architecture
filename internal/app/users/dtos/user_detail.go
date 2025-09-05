package dtos

import (
	"net/url"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/entities"
)

type UserDetailRequest interface {
	ToUserDetailOption() entities.UserDetailOption
}

type UserDetailByContactValueRequest struct {
	ContactValue string `param:"contact_value"`
}

type UserDetailByIDRequest struct {
	ID int64 `param:"id"`
}

func (u UserDetailByContactValueRequest) ToUserDetailOption() entities.UserDetailOption {
	decodedEmail, _ := url.QueryUnescape(u.ContactValue)
	return entities.WithContactValue(decodedEmail)
}

func (u UserDetailByIDRequest) ToUserDetailOption() entities.UserDetailOption {
	return entities.WithID(u.ID)
}

func ToUserDTO(u entities.User) User {
	return User{
		ID:           u.ID,
		Name:         u.Name,
		ContactType:  u.ContactType,
		ContactValue: u.ContactValue,
		BirthDate:    u.BirthDate,
		Language:     u.Language,
		Password:     u.Password,
		Status:       u.Status,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
