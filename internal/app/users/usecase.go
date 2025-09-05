package users

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/dtos"
)

type Usecase interface {
	Login(ctx context.Context, request dtos.UserLoginRequest) (response dtos.UserLoginResponse, err error)
	SignUp(ctx context.Context, request dtos.SignUpRequest) error
	UserDetail(ctx context.Context, request dtos.UserDetailRequest) (userData dtos.User, err error)
	UserUpdate(ctx context.Context, request dtos.UserUpdateRequest) error
	TransitionUserStatus(ctx context.Context, request dtos.TransitionUserStatusRequest) error
}
