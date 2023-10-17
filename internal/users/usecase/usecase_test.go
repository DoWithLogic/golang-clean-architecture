package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	mocks "github.com/DoWithLogic/golang-clean-architecture/internal/users/mock"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/usecase"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createUserMatcher(user entities.Users) gomock.Matcher {
	return eqUserMatcher{
		users: user,
	}
}

type eqUserMatcher struct {
	users entities.Users
}

func (e eqUserMatcher) Matches(x interface{}) bool {
	arg, ok := x.(entities.Users)
	if !ok {
		return false
	}

	return arg.Fullname == e.users.Fullname &&
		arg.PhoneNumber == e.users.PhoneNumber &&
		arg.UserType == e.users.UserType &&
		arg.IsActive == e.users.IsActive
}

func (e eqUserMatcher) String() string {
	return fmt.Sprintf("%v", e.users.Fullname)
}

func Test_usecase_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockRepository(ctrl)
	uc := usecase.NewUseCase(
		repo,
		zerolog.NewZeroLog(ctx, os.Stdout),
		config.Config{
			Authentication: config.AuthenticationConfig{
				Key:       "DoWithLogic!@#",
				SecretKey: "s3cr#tK3y!@#v001",
				SaltKey:   "s4ltK3y!@#ddv001",
			},
		},
	)

	newUser := dtos.CreateUserRequest{
		FullName:    "fullname",
		PhoneNumber: "081236548974",
		Email:       "martinyonatann@testing.com",
		Password:    "testingPwd",
	}

	t.Run("positive_case_create_user", func(t *testing.T) {
		repo.EXPECT().IsUserExist(ctx, newUser.Email).Return(false)

		repo.EXPECT().
			SaveNewUser(ctx,
				createUserMatcher(
					entities.Users{
						Fullname:    newUser.FullName,
						PhoneNumber: newUser.PhoneNumber,
						UserType:    constant.UserTypeRegular,
						IsActive:    true,
					},
				)).
			Return(int64(1), nil)

		userID, httpCode, err := uc.Create(ctx, newUser)
		require.NoError(t, err)
		require.Equal(t, httpCode, http.StatusOK)
		require.NotNil(t, userID)
	})

	t.Run("negative_email_already_use", func(t *testing.T) {
		repo.EXPECT().IsUserExist(ctx, newUser.Email).Return(true)

		userID, httpCode, err := uc.Create(ctx, newUser)
		require.EqualError(t, apperror.ErrEmailAlreadyExist, err.Error())
		require.Equal(t, httpCode, http.StatusConflict)
		require.Equal(t, userID, int64(0))
	})

	t.Run("negative_case_create_user_error_repo", func(t *testing.T) {
		repo.EXPECT().IsUserExist(ctx, newUser.Email).Return(false)

		repo.EXPECT().
			SaveNewUser(ctx,
				createUserMatcher(
					entities.Users{
						Fullname:    "fullname",
						PhoneNumber: "081236548974",
						UserType:    constant.UserTypeRegular,
						IsActive:    true,
					},
				)).
			Return(int64(0), sql.ErrNoRows)

		userID, httpCode, err := uc.Create(ctx, newUser)
		require.Error(t, err)
		require.EqualError(t, sql.ErrNoRows, err.Error())
		require.Equal(t, httpCode, http.StatusInternalServerError)
		require.Equal(t, userID, int64(0))
	})
}

func Test_usecase_UpdateUserStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockRepository(ctrl)
	uc := usecase.NewUseCase(
		repo,
		zerolog.NewZeroLog(ctx, os.Stdout),
		config.Config{Authentication: config.AuthenticationConfig{Key: "secret-key"}},
	)

	args := dtos.UpdateUserStatusRequest{
		UserID:   1,
		Status:   constant.UserActive,
		UpdateBy: "martin@test.com",
	}

	t.Run("positive_case_UpdateUserStatus", func(t *testing.T) {
		repo.EXPECT().
			GetUserByID(ctx, args.UserID, gomock.Any()).
			Return(entities.Users{UserID: 1, IsActive: true}, nil)

		repo.EXPECT().
			UpdateUserStatusByID(ctx, gomock.Any()).
			Return(nil)

		err := uc.UpdateStatus(ctx, args)
		require.NoError(t, err)
	})

	t.Run("negative_case_UpdateUserStatus_GetUserByID_err", func(t *testing.T) {
		repo.EXPECT().
			GetUserByID(ctx, args.UserID, gomock.Any()).
			Return(entities.Users{}, errors.New("something errors"))

		err := uc.UpdateStatus(ctx, args)
		require.Error(t, err)
	})

	t.Run("negative_case_UpdateUserStatus_err", func(t *testing.T) {
		repo.EXPECT().
			GetUserByID(ctx, args.UserID, gomock.Any()).
			Return(entities.Users{UserID: 1, IsActive: true}, nil)

		repo.EXPECT().
			UpdateUserStatusByID(ctx, gomock.Any()).
			Return(errors.New("there was error"))

		err := uc.UpdateStatus(ctx, args)
		require.Error(t, err)
	})

}

func Test_usecase_Detail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockRepository(ctrl)
	uc := usecase.NewUseCase(repo, zerolog.NewZeroLog(ctx, os.Stdout), config.Config{})

	var id int64 = 1

	returnedDetail := entities.Users{
		UserID:      id,
		Email:       "test@test.com",
		Fullname:    "test",
		PhoneNumber: "123456789012",
		UserType:    constant.UserTypePremium,
		IsActive:    true,
		CreatedAt:   time.Now(),
		CreatedBy:   "SYSTEM",
	}

	t.Run("detail_positive", func(t *testing.T) {
		repo.EXPECT().GetUserByID(ctx, id).Return(returnedDetail, nil)

		detail, httpCode, err := uc.Detail(ctx, id)
		require.NoError(t, err)
		require.Equal(t, httpCode, http.StatusOK)
		require.Equal(t, detail, entities.NewUserDetail(returnedDetail))
	})

	t.Run("detail_negative_failed_query_detail", func(t *testing.T) {
		repo.EXPECT().GetUserByID(ctx, id).Return(entities.Users{}, sql.ErrNoRows)

		detail, httpCode, err := uc.Detail(ctx, id)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Equal(t, httpCode, http.StatusInternalServerError)
		require.Equal(t, detail, dtos.UserDetailResponse{})
	})

}

func Test_usecase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		password = "testing"
		email    = "martin@test.com"

		config = config.Config{
			Authentication: config.AuthenticationConfig{
				Key:       "DoWithLogic!@#",
				SecretKey: "s3cr#tK3y!@#v001",
				SaltKey:   "s4ltK3y!@#ddv001",
			},
		}
	)

	ctx := context.Background()
	repo := mocks.NewMockRepository(ctrl)
	uc := usecase.NewUseCase(repo, zerolog.NewZeroLog(ctx, os.Stdout), config)

	returnedUser := entities.Users{
		UserID:   1,
		Email:    email,
		Password: utils.Encrypt(password, config),
	}

	t.Run("login_positive", func(t *testing.T) {
		repo.EXPECT().GetUserByEmail(ctx, email).Return(returnedUser, nil)

		authData, code, err := uc.Login(ctx, dtos.UserLoginRequest{Email: email, Password: password})
		require.NoError(t, err)
		require.Equal(t, code, http.StatusOK)
		require.NotNil(t, authData)

	})

	t.Run("login_negative_invalid_password", func(t *testing.T) {
		repo.EXPECT().GetUserByEmail(ctx, email).Return(returnedUser, nil)

		authData, code, err := uc.Login(ctx, dtos.UserLoginRequest{Email: email, Password: "testingpwd"})
		require.EqualError(t, apperror.ErrInvalidPassword, err.Error())
		require.Equal(t, code, http.StatusUnauthorized)
		require.Equal(t, authData, dtos.UserLoginResponse{})

	})

	t.Run("login_negative_failed_query_email", func(t *testing.T) {
		repo.EXPECT().GetUserByEmail(ctx, email).Return(entities.Users{}, sql.ErrNoRows)

		authData, code, err := uc.Login(ctx, dtos.UserLoginRequest{Email: email, Password: password})
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Equal(t, code, http.StatusInternalServerError)
		require.Equal(t, authData, dtos.UserLoginResponse{})

	})

}

func Test_usecase_PartialUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		request = dtos.UpdateUserRequest{
			Fullname: "update name",
			UserID:   1,
		}

		config = config.Config{
			Authentication: config.AuthenticationConfig{
				Key:       "DoWithLogic!@#",
				SecretKey: "s3cr#tK3y!@#v001",
				SaltKey:   "s4ltK3y!@#ddv001",
			},
		}
	)

	ctx := context.Background()
	repo := mocks.NewMockRepository(ctrl)
	uc := usecase.NewUseCase(repo, zerolog.NewZeroLog(ctx, os.Stdout), config)

	repo.EXPECT().Atomic(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	err := uc.PartialUpdate(context.Background(), request)
	require.NoError(t, err)
}
