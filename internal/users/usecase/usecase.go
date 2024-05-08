package usecase

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/middleware"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_crypto"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/golang-jwt/jwt"
)

type (
	usecase struct {
		repo users.Repository
		cfg  config.Config
	}
)

func NewUseCase(repo users.Repository, cfg config.Config) users.Usecase {
	return &usecase{repo, cfg}
}

func (uc *usecase) Login(ctx context.Context, request dtos.UserLoginRequest) (response dtos.UserLoginResponse, err error) {
	dataLogin, err := uc.repo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return response, apperror.InternalServerError(err)
	}

	if !strings.EqualFold(dataLogin.Password, app_crypto.NewCrypto(uc.cfg.Authentication.Key).EncodeSHA256(request.Password)) {
		return response, apperror.Unauthorized(apperror.ErrInvalidPassword)
	}

	claims := middleware.PayloadToken{
		Data: &middleware.Data{
			UserID: dataLogin.UserID,
			Email:  dataLogin.Email,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}

	// Calculate the expiration time in seconds
	expiresIn := claims.ExpiresAt - time.Now().Unix()

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString([]byte(uc.cfg.JWT.Key))
	if err != nil {
		return response, apperror.InternalServerError(err)
	}

	return dtos.UserLoginResponse{AccessToken: tokenString, ExpiredAt: expiresIn}, nil
}

func (uc *usecase) Create(ctx context.Context, payload dtos.CreateUserRequest) (userID int64, err error) {
	if exist := uc.repo.IsUserExist(ctx, payload.Email); exist {
		return userID, apperror.Conflict(apperror.ErrEmailAlreadyExist)
	}

	userID, err = uc.repo.SaveNewUser(ctx, entities.NewCreateUser(payload, uc.cfg))
	if err != nil {
		return userID, apperror.InternalServerError(err)
	}

	return userID, nil
}

func (uc *usecase) PartialUpdate(ctx context.Context, data dtos.UpdateUserRequest) error {
	return uc.repo.Atomic(ctx, &sql.TxOptions{}, func(tx users.Repository) error {
		opt := entities.LockingOpt{
			PessimisticLocking: true,
		}

		if _, err := tx.GetUserByID(ctx, data.UserID, opt); err != nil {
			return err
		}

		return tx.UpdateUserByID(ctx, entities.NewUpdateUser(data))
	})
}

func (uc *usecase) UpdateStatus(ctx context.Context, req dtos.UpdateUserStatusRequest) error {
	if _, err := uc.repo.GetUserByID(ctx, req.UserID, entities.LockingOpt{}); err != nil {
		return err
	}

	return uc.repo.UpdateUserStatusByID(ctx, entities.NewUpdateUserStatus(req))
}

func (uc *usecase) Detail(ctx context.Context, id int64) (detail dtos.UserDetailResponse, err error) {
	userDetail, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		return detail, apperror.InternalServerError(err)
	}

	return entities.NewUserDetail(userDetail), nil
}
