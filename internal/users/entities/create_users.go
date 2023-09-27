package entities

import "time"

type CreateUser struct {
	FullName    string
	PhoneNumber string
	UserType    string
	IsActive    bool
	CreatedAt   time.Time
	CreatedBy   string
}
