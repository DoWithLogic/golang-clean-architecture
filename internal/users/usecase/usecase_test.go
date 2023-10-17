package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	mocks "github.com/DoWithLogic/golang-clean-architecture/internal/users/mock"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/usecase"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
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
		config.Config{Authentication: config.AuthenticationConfig{Key: "DoWithLogic!@#", SecretKey: "s3cr#tK3y!@#"}},
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

	t.Run("negative_case_create_user_error_repo", func(t *testing.T) {
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
			Return(int64(0), errors.New("something errors"))

		userID, httpCode, err := uc.Create(ctx, newUser)
		require.Error(t, err)
		require.Equal(t, httpCode, http.StatusInternalServerError)
		require.Equal(t, userID, 0)
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

	args := entities.UpdateUserStatus{
		UserID:   1,
		IsActive: true,
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
