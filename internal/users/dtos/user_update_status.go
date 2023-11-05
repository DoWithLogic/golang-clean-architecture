package dtos

import (
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
	"github.com/invopop/validation"
)

type UpdateUserStatusRequest struct {
	UserID   int64  `json:"-"`
	Status   string `json:"status"`
	UpdateBy string `json:"-"`
}

func (ussp UpdateUserStatusRequest) Validate() error {
	if ussp.Status != constant.UserInactive && ussp.Status != constant.UserActive {
		return apperror.ErrStatusValue
	}

	return validation.ValidateStruct(&ussp,
		validation.Field(&ussp.UserID, validation.NotNil),
		validation.Field(&ussp.UpdateBy, validation.NotNil),
	)
}
