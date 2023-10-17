package usecase

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/repository"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/middleware"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/utils"
	"github.com/dgrijalva/jwt-go"
)

type (
	Usecase interface {
		Login(ctx context.Context, request dtos.UserLoginRequest) (response dtos.UserLoginResponse, httpCode int, err error)
		Create(ctx context.Context, payload dtos.CreateUserRequest) (userID int64, httpCode int, err error)
		PartialUpdate(ctx context.Context, data dtos.UpdateUserRequest) error
		UpdateStatus(ctx context.Context, req dtos.UpdateUserStatusRequest) error
		Detail(ctx context.Context, id int64) (detail dtos.UserDetailResponse, httpCode int, err error)
	}

	usecase struct {
		repo repository.Repository
		log  *zerolog.Logger
		cfg  config.Config
	}
)

func NewUseCase(repo repository.Repository, log *zerolog.Logger, cfg config.Config) Usecase {
	return &usecase{repo, log, cfg}
}

func (uc *usecase) Login(ctx context.Context, request dtos.UserLoginRequest) (response dtos.UserLoginResponse, httpCode int, err error) {
	dataLogin, err := uc.repo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	if !strings.EqualFold(utils.Decrypt(dataLogin.Password, uc.cfg), request.Password) {
		return response, http.StatusUnauthorized, apperror.ErrInvalidPassword
	}

	identityData := middleware.CustomClaims{
		UserID: dataLogin.UserID,
		Email:  dataLogin.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	token, err := middleware.GenerateJWT(identityData, uc.cfg.Authentication.Key)
	if err != nil {
		return response, http.StatusInternalServerError, apperror.ErrFailedGenerateJWT
	}

	response = dtos.UserLoginResponse{
		AccessToken: token,
		ExpiredAt:   utils.UnixToDuration(identityData.ExpiresAt),
	}

	return response, http.StatusOK, nil
}

func (uc *usecase) Create(ctx context.Context, payload dtos.CreateUserRequest) (userID int64, httpCode int, err error) {
	if exist := uc.repo.IsUserExist(ctx, payload.Email); exist {
		return userID, http.StatusConflict, apperror.ErrEmailAlreadyExist
	}

	userID, err = uc.repo.SaveNewUser(ctx, entities.NewCreateUser(payload, uc.cfg))
	if err != nil {
		uc.log.Z().Err(err).Msg("[usecase]CreateUser.SaveNewUser")

		return userID, http.StatusInternalServerError, err
	}

	return userID, http.StatusOK, nil
}

func (uc *usecase) PartialUpdate(ctx context.Context, data dtos.UpdateUserRequest) error {
	return uc.repo.Atomic(ctx, &sql.TxOptions{}, func(tx repository.Repository) error {
		opt := entities.LockingOpt{
			PessimisticLocking: true,
		}
		_, err := tx.GetUserByID(ctx, data.UserID, opt)
		if err != nil {
			uc.log.Z().Err(err).Msg("[usecase]UpdateUser.GetUserByID")

			return err
		}

		err = tx.UpdateUserByID(ctx, entities.NewUpdateUsers(data))
		if err != nil {
			uc.log.Z().Err(err).Msg("[usecase]UpdateUser.UpdateUserByID")

			return err
		}

		return nil
	})
}

func (uc *usecase) UpdateStatus(ctx context.Context, req dtos.UpdateUserStatusRequest) error {
	_, err := uc.repo.GetUserByID(ctx, req.UserID, entities.LockingOpt{})
	if err != nil {
		uc.log.Z().Err(err).Msg("[usecase]UpdateUserStatus.GetUserByID")

		return err
	}

	if err := uc.repo.UpdateUserStatusByID(ctx, entities.NewUpdateUserStatus(req)); err != nil {
		uc.log.Z().Err(err).Msg("[usecase]UpdateUserStatus.UpdateUserStatusByID")

		return err
	}

	return nil
}

func (uc *usecase) Detail(ctx context.Context, id int64) (detail dtos.UserDetailResponse, httpCode int, err error) {
	userDetail, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		return detail, http.StatusInternalServerError, err
	}

	return entities.NewUserDetail(userDetail), http.StatusOK, nil
}
