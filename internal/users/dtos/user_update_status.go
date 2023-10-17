package dtos

import "github.com/invopop/validation"

type UpdateUserStatusRequest struct {
	UserID   int64  `json:"-"`
	Status   int    `json:"status"`
	UpdateBy string `json:"-"`
}

func (ussp UpdateUserStatusRequest) Validate() error {
	return validation.ValidateStruct(&ussp,
		validation.Field(&ussp.UserID, validation.NotNil),
		validation.Field(&ussp.UpdateBy, validation.NotNil),
		validation.Field(&ussp.Status, validation.Required),
	)
}
