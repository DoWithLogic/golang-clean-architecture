package entities

import "time"

type UpdateUserStatus struct {
	UserID    int64
	IsActive  bool
	UpdatedAt time.Time
	UpdatedBy string
}

func NewUpdateUserStatus(payload UpdateUserStatus) UpdateUserStatus {
	return UpdateUserStatus{
		UserID:    payload.UserID,
		IsActive:  payload.IsActive,
		UpdatedAt: time.Now(),
		UpdatedBy: "martin",
	}
}
