package entities

import "time"

type UpdateUsers struct {
	UserID      int64
	Fullname    string
	PhoneNumber string
	UserType    string
	IsActive    bool
	UpdatedAt   string
	UpdatedBy   string
}

func NewUpdateUsers(data UpdateUsers) *UpdateUsers {
	return &UpdateUsers{
		Fullname:    data.Fullname,
		PhoneNumber: data.PhoneNumber,
		UserType:    UserTypeRegular,
		IsActive:    true,
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedBy:   "martin",
	}
}
