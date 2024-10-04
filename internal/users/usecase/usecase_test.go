package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	mocks "github.com/DoWithLogic/golang-clean-architecture/internal/users/mock"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/usecase"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_crypto"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_jwt"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	KeyUnitTest = "DoWithLogic!@#"
)

func createUserMatcher(user entities.User) gomock.Matcher {
	return eqUserMatcher{
		users: user,
	}
}

type eqUserMatcher struct {
	users entities.User
}

func (e eqUserMatcher) Matches(x interface{}) bool {
	arg, ok := x.(entities.User)
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

	crypto := app_crypto.NewCrypto(KeyUnitTest)
	appJwt := app_jwt.NewJWT(config.JWTConfig{Key: KeyUnitTest, Expired: 60, Label: "XXXX"})
	ctx := context.Background()
	repo := mocks.NewMockRepository(ctrl)
	uc := usecase.NewUseCase(
		repo,
		appJwt,
		crypto,
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
					entities.User{
						Fullname:    newUser.FullName,
						PhoneNumber: newUser.PhoneNumber,
						UserType:    constant.UserTypeRegular,
						IsActive:    true,
					},
				)).
			Return(int64(1), nil)

		userID, err := uc.Create(ctx, newUser)
		require.NoError(t, err)
		require.NotNil(t, userID)
	})

	t.Run("negative_email_already_use", func(t *testing.T) {
		repo.EXPECT().IsUserExist(ctx, newUser.Email).Return(true)

		userID, err := uc.Create(ctx, newUser)
		require.EqualError(t, apperror.ErrEmailAlreadyExist, err.Error())
		require.Equal(t, userID, int64(0))
	})

	t.Run("negative_case_create_user_error_repo", func(t *testing.T) {
		repo.EXPECT().IsUserExist(ctx, newUser.Email).Return(false)

		repo.EXPECT().
			SaveNewUser(ctx,
				createUserMatcher(
					entities.User{
						Fullname:    "fullname",
						PhoneNumber: "081236548974",
						UserType:    constant.UserTypeRegular,
						IsActive:    true,
					},
				)).
			Return(int64(0), sql.ErrNoRows)

		userID, err := uc.Create(ctx, newUser)
		require.Error(t, err)
		require.EqualError(t, sql.ErrNoRows, err.Error())
		require.Equal(t, userID, int64(0))
	})
}

func Test_usecase_UpdateUserStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	crypto := app_crypto.NewCrypto(KeyUnitTest)
	appJwt := app_jwt.NewJWT(config.JWTConfig{Key: KeyUnitTest, Expired: 60, Label: "XXXX"})
	ctx := context.Background()
	repo := mocks.NewMockRepository(ctrl)
	uc := usecase.NewUseCase(
		repo,
		appJwt,
		crypto,
	)

	args := dtos.UpdateUserStatusRequest{
		UserID: 1,
		UpdateUserStatus: dtos.UpdateUserStatus{
			Status: constant.UserActive,
		},
	}

	t.Run("positive_case_UpdateUserStatus", func(t *testing.T) {
		repo.EXPECT().
			GetUserByID(ctx, args.UserID, gomock.Any()).
			Return(entities.User{UserID: 1, IsActive: true}, nil)

		repo.EXPECT().
			UpdateUserStatusByID(ctx, gomock.Any()).
			Return(nil)

		err := uc.UpdateStatus(ctx, args)
		require.NoError(t, err)
	})

	t.Run("negative_case_UpdateUserStatus_GetUserByID_err", func(t *testing.T) {
		repo.EXPECT().
			GetUserByID(ctx, args.UserID, gomock.Any()).
			Return(entities.User{}, errors.New("something errors"))

		err := uc.UpdateStatus(ctx, args)
		require.Error(t, err)
	})

	t.Run("negative_case_UpdateUserStatus_err", func(t *testing.T) {
		repo.EXPECT().
			GetUserByID(ctx, args.UserID, gomock.Any()).
			Return(entities.User{UserID: 1, IsActive: true}, nil)

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

	crypto := app_crypto.NewCrypto(KeyUnitTest)
	appJwt := app_jwt.NewJWT(config.JWTConfig{Key: KeyUnitTest, Expired: 60, Label: "XXXX"})
	ctx := context.Background()
	repo := mocks.NewMockRepository(ctrl)
	uc := usecase.NewUseCase(
		repo,
		appJwt,
		crypto,
	)

	var id int64 = 1

	returnedDetail := entities.User{
		UserID:      id,
		Email:       "test@test.com",
		Fullname:    "test",
		PhoneNumber: "123456789012",
		UserType:    constant.UserTypePremium,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	t.Run("detail_positive", func(t *testing.T) {
		repo.EXPECT().GetUserByID(ctx, id).Return(returnedDetail, nil)

		detail, err := uc.Detail(ctx, id)
		require.NoError(t, err)
		require.Equal(t, detail, entities.NewUserDetail(returnedDetail))
	})

	t.Run("detail_negative_failed_query_detail", func(t *testing.T) {
		repo.EXPECT().GetUserByID(ctx, id).Return(entities.User{}, sql.ErrNoRows)

		detail, err := uc.Detail(ctx, id)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Equal(t, detail, dtos.UserDetailResponse{})
	})

}

func Test_usecase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		password = "testing"
		email    = "martin@test.com"
	)

	crypto := app_crypto.NewCrypto(KeyUnitTest)
	appJwt := app_jwt.NewJWT(config.JWTConfig{Key: KeyUnitTest, Expired: 60, Label: "XXXX"})
	ctx := context.Background()
	repo := mocks.NewMockRepository(ctrl)
	uc := usecase.NewUseCase(
		repo,
		appJwt,
		crypto,
	)

	returnedUser := entities.User{
		UserID:   1,
		Email:    email,
		Password: crypto.EncodeSHA256(password),
	}

	t.Run("login_positive", func(t *testing.T) {
		repo.EXPECT().GetUserByEmail(ctx, email).Return(returnedUser, nil)

		authData, err := uc.Login(ctx, dtos.UserLoginRequest{Email: email, Password: password})
		require.NoError(t, err)
		require.NotNil(t, authData)

	})

	t.Run("login_negative_invalid_password", func(t *testing.T) {
		repo.EXPECT().GetUserByEmail(ctx, email).Return(returnedUser, nil)

		authData, err := uc.Login(ctx, dtos.UserLoginRequest{Email: email, Password: "testingpwd"})
		require.EqualError(t, apperror.ErrInvalidPassword, err.Error())
		require.Equal(t, authData, dtos.UserLoginResponse{})

	})

	t.Run("login_negative_failed_query_email", func(t *testing.T) {
		repo.EXPECT().GetUserByEmail(ctx, email).Return(entities.User{}, sql.ErrNoRows)

		authData, err := uc.Login(ctx, dtos.UserLoginRequest{Email: email, Password: password})
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Equal(t, authData, dtos.UserLoginResponse{})

	})

}

func Test_usecase_PartialUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		request = dtos.UpdateUserRequest{
			UpdateUser: dtos.UpdateUser{
				Fullname:    faker.Name(),
				PhoneNumber: faker.Phonenumber(),
			},
			UserID: 1,
		}
	)

	crypto := app_crypto.NewCrypto(KeyUnitTest)
	appJwt := app_jwt.NewJWT(config.JWTConfig{Key: KeyUnitTest, Expired: 60, Label: "XXXX"})
	repo := mocks.NewMockRepository(ctrl)
	uc := usecase.NewUseCase(
		repo,
		appJwt,
		crypto,
	)

	repo.EXPECT().Atomic(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	err := uc.PartialUpdate(context.Background(), request)
	require.NoError(t, err)
}
