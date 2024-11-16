package usecase

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_crypto"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_jwt"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability/instrumentation"
	"github.com/golang-jwt/jwt"
)

type (
	usecase struct {
		repo   users.Repository
		appJwt *app_jwt.JWT
		crypto *app_crypto.Crypto
	}
)

func NewUseCase(repo users.Repository, appJwt *app_jwt.JWT, crypto *app_crypto.Crypto) users.Usecase {
	return &usecase{repo, appJwt, crypto}
}

func (uc *usecase) Login(ctx context.Context, request dtos.UserLoginRequest) (response dtos.UserLoginResponse, err error) {
	ctx, span := instrumentation.NewTraceSpan(ctx, "LoginUC")
	defer span.End()

	dataLogin, err := uc.repo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return response, apperror.InternalServerError(err)
	}

	if !strings.EqualFold(dataLogin.Password, uc.crypto.EncodeSHA256(request.Password)) {
		return response, apperror.Unauthorized(apperror.ErrInvalidPassword)
	}

	claims := app_jwt.PayloadToken{
		Data: &app_jwt.Data{
			UserID: dataLogin.UserID,
			Email:  dataLogin.Email,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}

	// Generate encoded token and send it as response.
	jwtToken, err := uc.appJwt.GenerateToken(ctx, claims)
	if err != nil {
		return response, apperror.InternalServerError(err)
	}

	return dtos.UserLoginResponse{AccessToken: jwtToken, ExpiredAt: claims.ExpiresAt - time.Now().Unix()}, nil
}

func (uc *usecase) Create(ctx context.Context, payload dtos.CreateUserRequest) (userID int64, err error) {
	ctx, span := instrumentation.NewTraceSpan(ctx, "CreateUC")
	defer span.End()

	if exist := uc.repo.IsUserExist(ctx, payload.Email); exist {
		return userID, apperror.Conflict(apperror.ErrEmailAlreadyExist)
	}

	payload.Password = uc.crypto.EncodeSHA256(payload.Password)
	userID, err = uc.repo.SaveNewUser(ctx, entities.NewCreateUser(payload))
	if err != nil {
		return userID, apperror.InternalServerError(err)
	}

	return userID, nil
}

func (uc *usecase) PartialUpdate(ctx context.Context, data dtos.UpdateUserRequest) error {
	ctx, span := instrumentation.NewTraceSpan(ctx, "PartialUpdateUC")
	defer span.End()

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
	ctx, span := instrumentation.NewTraceSpan(ctx, "UpdateStatusUC")
	defer span.End()

	if _, err := uc.repo.GetUserByID(ctx, req.UserID, entities.LockingOpt{}); err != nil {
		return err
	}

	return uc.repo.UpdateUserStatusByID(ctx, entities.NewUpdateUserStatus(req))
}

func (uc *usecase) Detail(ctx context.Context, id int64) (detail dtos.UserDetailResponse, err error) {
	ctx, span := instrumentation.NewTraceSpan(ctx, "DetailUC")
	defer span.End()

	userDetail, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		return detail, apperror.InternalServerError(err)
	}

	return entities.NewUserDetail(userDetail), nil
}
