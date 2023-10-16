package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	mocks "github.com/DoWithLogic/golang-clean-architecture/internal/users/mock"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/usecase"
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
	)

	newUser := entities.CreateUser{
		FullName:    "fullname",
		PhoneNumber: "081236548974",
		UserType:    entities.UserTypePremium,
		IsActive:    true,
	}

	t.Run("positive_case_create_user", func(t *testing.T) {

		repo.EXPECT().
			SaveNewUser(ctx,
				createUserMatcher(
					entities.Users{
						Fullname:    "fullname",
						PhoneNumber: "081236548974",
						UserType:    entities.UserTypeRegular,
						IsActive:    true,
					},
				)).
			Return(int64(1), nil)

		userID, err := uc.CreateUser(ctx, newUser)
		require.NoError(t, err)
		require.NotNil(t, userID)
	})

	t.Run("negative_case_create_user_error_repo", func(t *testing.T) {
		repo.EXPECT().
			SaveNewUser(ctx,
				createUserMatcher(
					entities.Users{
						Fullname:    "fullname",
						PhoneNumber: "081236548974",
						UserType:    entities.UserTypeRegular,
						IsActive:    true,
					},
				)).
			Return(int64(0), errors.New("something errors"))

		userID, err := uc.CreateUser(ctx, newUser)
		require.Error(t, err)
		require.Equal(t, userID, int64(0))
	})
}

func Test_usecase_UpdateUserStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := mocks.NewMockRepository(ctrl)
	uc := usecase.NewUseCase(repo, zerolog.NewZeroLog(ctx, os.Stdout))

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

		err := uc.UpdateUserStatus(ctx, args)
		require.NoError(t, err)
	})

	t.Run("negative_case_UpdateUserStatus_GetUserByID_err", func(t *testing.T) {
		repo.EXPECT().
			GetUserByID(ctx, args.UserID, gomock.Any()).
			Return(entities.Users{}, errors.New("something errors"))

		err := uc.UpdateUserStatus(ctx, args)
		require.Error(t, err)
	})

	t.Run("negative_case_UpdateUserStatus_err", func(t *testing.T) {
		repo.EXPECT().
			GetUserByID(ctx, args.UserID, gomock.Any()).
			Return(entities.Users{UserID: 1, IsActive: true}, nil)

		repo.EXPECT().
			UpdateUserStatusByID(ctx, gomock.Any()).
			Return(errors.New("there was error"))

		err := uc.UpdateUserStatus(ctx, args)
		require.Error(t, err)
	})

}
