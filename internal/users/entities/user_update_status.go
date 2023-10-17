package entities

import (
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
)

type UpdateUserStatus struct {
	UserID    int64
	IsActive  bool
	UpdatedAt time.Time
	UpdatedBy string
}

func NewUpdateUserStatus(req dtos.UpdateUserStatusRequest) UpdateUserStatus {
	return UpdateUserStatus{
		UserID:    req.UserID,
		IsActive:  constant.MapStatus[req.Status],
		UpdatedAt: time.Now(),
		UpdatedBy: req.UpdateBy,
	}
}
