package users

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
)

type Usecase interface {
	Login(ctx context.Context, request dtos.UserLoginRequest) (response dtos.UserLoginResponse, httpCode int, err error)
	Create(ctx context.Context, payload dtos.CreateUserRequest) (userID int64, httpCode int, err error)
	PartialUpdate(ctx context.Context, data dtos.UpdateUserRequest) error
	UpdateStatus(ctx context.Context, req dtos.UpdateUserStatusRequest) error
	Detail(ctx context.Context, id int64) (detail dtos.UserDetailResponse, httpCode int, err error)
}
