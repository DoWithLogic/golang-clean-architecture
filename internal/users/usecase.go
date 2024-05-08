package users

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
)

type Usecase interface {
	Login(ctx context.Context, request dtos.UserLoginRequest) (response dtos.UserLoginResponse, err error)
	Create(ctx context.Context, payload dtos.CreateUserRequest) (userID int64, err error)
	PartialUpdate(ctx context.Context, data dtos.UpdateUserRequest) error
	UpdateStatus(ctx context.Context, req dtos.UpdateUserStatusRequest) error
	Detail(ctx context.Context, id int64) (detail dtos.UserDetailResponse, err error)
}
