package entities

import "time"

type UpdateUsers struct {
	UserID      int64
	Fullname    string
	PhoneNumber string
	UserType    string
	UpdatedAt   time.Time
	UpdatedBy   string
}

func NewUpdateUsers(data UpdateUsers) UpdateUsers {
	return UpdateUsers{
		UserID:      data.UserID,
		Fullname:    data.Fullname,
		PhoneNumber: data.PhoneNumber,
		UserType:    data.UserType,
		UpdatedAt:   time.Now(),
		UpdatedBy:   "martin",
	}
}
