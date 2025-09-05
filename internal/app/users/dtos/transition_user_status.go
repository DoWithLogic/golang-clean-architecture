package dtos

import (
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
	"github.com/invopop/validation"
)

type TransitionUserStatus struct {
	Status types.USER_STATUS `json:"status"`
}

type TransitionUserStatusRequest struct {
	ID int64 `param:"id"`
	TransitionUserStatus
}

func (u TransitionUserStatusRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ID, validation.Required),
		validation.Field(&u.Status, validation.Required, validation.In(types.ACTIVE, types.PENDING, types.REJECT, types.BANNED)),
	)
}

func (t TransitionUserStatusRequest) ToUpdateStatusEntity() *entities.UpdateUser {
	return &entities.UpdateUser{
		ID:        t.ID,
		Status:    &t.Status,
		UpdatedAt: time.Now(),
	}
}
