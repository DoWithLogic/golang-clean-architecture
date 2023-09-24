package entities

import "time"

type UpdateUsers struct {
	UserID      int64
	Fullname    string
	PhoneNumber string
	UserType    string
	UpdatedAt   string
	UpdatedBy   string
}

func NewUpdateUsers(data UpdateUsers) *UpdateUsers {
	return &UpdateUsers{
		UserID:      data.UserID,
		Fullname:    data.Fullname,
		PhoneNumber: data.PhoneNumber,
		UserType:    UserTypeRegular,
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedBy:   "martin",
	}
}
