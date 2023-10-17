package dtos

import "time"

type (
	UserDetailResponse struct {
		UserID      int64     `json:"id"`
		Email       string    `json:"email"`
		Fullname    string    `json:"fullname"`
		PhoneNumber string    `json:"phone_number"`
		UserType    string    `json:"user_type"`
		IsActive    bool      `json:"is_active"`
		CreatedAt   time.Time `json:"created_at"`
		CreatedBy   string    `json:"created_by"`
	}
)
