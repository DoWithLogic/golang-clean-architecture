package entities

import "time"

type CreateUser struct {
	FUllName    string
	PhoneNumber string
	UserType    string
	IsActive    bool
	CreatedAt   time.Time
	CreatedBy   string
}

func NewCreateUser(data CreateUser) CreateUser {
	return CreateUser{
		FUllName:    data.FUllName,
		PhoneNumber: data.PhoneNumber,
		UserType:    UserTypeRegular,
		IsActive:    true,
		CreatedAt:   time.Now(),
		CreatedBy:   "martin",
	}
}
