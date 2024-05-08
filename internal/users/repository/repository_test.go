package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/repository"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/constant"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasource"
	"github.com/go-faker/faker/v4"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

var cfg config.Config
var db *sqlx.DB

func init() {
	cfg = lo.Must(config.LoadConfigPath("../../../config/config-local"))
	db = lo.Must(datasource.NewDatabase(cfg.Database))

}

func Test_repository_SaveNewUser(t *testing.T) {
	repo := repository.NewRepository(db)

	t.Run("positive_SaveNewUser", func(t *testing.T) {
		newUserData := entities.User{
			Email:       faker.Email(),
			Password:    faker.Password(),
			Fullname:    faker.Name(),
			PhoneNumber: faker.Phonenumber(),
			UserType:    constant.UserTypePremium,
			CreatedAt:   time.Now(),
		}

		userID, err := repo.SaveNewUser(context.Background(), newUserData)
		require.NoError(t, err)
		require.NotEmpty(t, userID)

		userData, err := repo.GetUserByID(context.Background(), userID)
		require.NoError(t, err)
		require.Equal(t, userData.Email, newUserData.Email)
		require.Equal(t, userData.Fullname, newUserData.Fullname)
		require.Equal(t, userData.PhoneNumber, newUserData.PhoneNumber)
		require.Equal(t, userData.UserType, newUserData.UserType)
	})

}

func Test_repository_UpdateUserByID(t *testing.T) {
	repo := repository.NewRepository(db)

	t.Run("positive_UpdateUserByID", func(t *testing.T) {
		newUserData := entities.User{
			Email:       faker.Email(),
			Password:    faker.Password(),
			Fullname:    faker.Name(),
			PhoneNumber: faker.Phonenumber(),
			UserType:    constant.UserTypePremium,
			CreatedAt:   time.Now(),
		}

		userID, err := repo.SaveNewUser(context.Background(), newUserData)
		require.NoError(t, err)
		require.NotEmpty(t, userID)

		userData, err := repo.GetUserByID(context.Background(), userID)
		require.NoError(t, err)
		require.Equal(t, userData.Email, newUserData.Email)
		require.Equal(t, userData.Fullname, newUserData.Fullname)
		require.Equal(t, userData.PhoneNumber, newUserData.PhoneNumber)
		require.Equal(t, userData.UserType, newUserData.UserType)

		updateUserData := entities.UpdateUser{
			UserID:      userID,
			Fullname:    "updated name",
			PhoneNumber: "081212121313",
			UserType:    constant.UserTypeRegular,
			UpdatedAt:   time.Now(),
		}

		err = repo.UpdateUserByID(context.Background(), updateUserData)
		require.NoError(t, err)

		userDataAfterUpdate, err := repo.GetUserByID(context.Background(), userID)
		require.NoError(t, err)
		require.Equal(t, userDataAfterUpdate.Fullname, updateUserData.Fullname)
		require.Equal(t, userDataAfterUpdate.PhoneNumber, updateUserData.PhoneNumber)
		require.Equal(t, userDataAfterUpdate.UserType, updateUserData.UserType)

	})

}

func Test_repository_UpdateUserStatusByID(t *testing.T) {
	repo := repository.NewRepository(db)

	t.Run("positive_UpdateUserStatusByID", func(t *testing.T) {
		newUserData := entities.User{
			Email:       faker.Email(),
			Password:    faker.Password(),
			Fullname:    faker.Name(),
			PhoneNumber: faker.Phonenumber(),
			UserType:    constant.UserTypePremium,
			CreatedAt:   time.Now(),
		}

		userID, err := repo.SaveNewUser(context.Background(), newUserData)
		require.NoError(t, err)
		require.NotEmpty(t, userID)

		exist := repo.IsUserExist(context.Background(), newUserData.Email)
		require.Equal(t, exist, true)

		updateUserStatusData := entities.UpdateUserStatus{
			UserID:    userID,
			IsActive:  true,
			UpdatedAt: time.Now(),
		}

		repo.Atomic(context.Background(), &sql.TxOptions{}, func(tx users.Repository) error {
			_, err := tx.GetUserByID(context.Background(), userID, entities.LockingOpt{PessimisticLocking: true})
			require.NoError(t, err)

			err = tx.UpdateUserStatusByID(context.Background(), updateUserStatusData)
			require.NoError(t, err)

			return nil
		})

		userDataAfterUpdateStatus, err := repo.GetUserByEmail(context.Background(), newUserData.Email)
		require.NoError(t, err)
		require.Equal(t, userDataAfterUpdateStatus.IsActive, updateUserStatusData.IsActive)
	})
}
