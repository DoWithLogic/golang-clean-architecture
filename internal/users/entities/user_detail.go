package entities

import "github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"

func NewUserDetail(data Users) dtos.UserDetailResponse {
	return dtos.UserDetailResponse{
		UserID:      data.UserID,
		Email:       data.Email,
		Fullname:    data.Fullname,
		PhoneNumber: data.PhoneNumber,
		UserType:    data.UserType,
		IsActive:    data.IsActive,
		CreatedAt:   data.CreatedAt,
		CreatedBy:   data.CreatedBy,
	}
}
