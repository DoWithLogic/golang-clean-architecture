package dtos

import (
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
	"github.com/invopop/validation"
)

type UpdateUserStatus struct {
	Status string `json:"status"`
}

type UpdateUserStatusRequest struct {
	UserID int64
	UpdateUserStatus
}

func (ussp UpdateUserStatus) Validate() error {
	return validation.ValidateStruct(&ussp,
		validation.Field(&ussp.Status, validation.In(constant.UserActive, constant.UserInactive)),
	)
}
