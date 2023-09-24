package repository

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/repository/repository_query"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/custom"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/database"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
)

type (
	Repository interface {
		SaveNewUser(context.Context, *entities.Users) (int64, error)
		UpdateUserByID(context.Context, *entities.UpdateUsers) error
		GetUserByID(context.Context, int64, entities.LockingOpt) (entities.Users, error)
	}

	repository struct {
		conn database.SQLTxConn
		log  *zerolog.Logger
	}
)

func NewRepository(conn database.SQLTxConn, log *zerolog.Logger) Repository {
	return &repository{conn, log}
}

func (repo *repository) SaveNewUser(ctx context.Context, user *entities.Users) (int64, error) {
	args := custom.Array{
		user.Fullname,
		user.PhoneNumber,
		user.IsActive,
		user.UserType,
		user.CreatedAt,
		user.CreatedBy,
	}

	var userID int64
	err := new(database.SQL).Exec(repo.conn.ExecContext(ctx, repository_query.InsertUsers, args...)).Scan(nil, &userID)
	if err != nil {
		repo.log.Z().Err(err).Msg("[repository]SaveNewUser.ExecContext")

		return userID, err
	}

	return userID, err
}

func (repo *repository) UpdateUserByID(ctx context.Context, user *entities.UpdateUsers) error {
	args := custom.Array{
		user.Fullname, user.Fullname,
		user.PhoneNumber, user.PhoneNumber,
		user.UserType, user.UserType,
		user.UpdatedAt,
		user.UpdatedBy,
		user.UserID,
	}

	err := new(database.SQL).Exec(repo.conn.ExecContext(ctx, repository_query.UpdateUsers, args...)).Scan(nil, nil)
	if err != nil {
		repo.log.Z().Err(err).Msg("[repository]UpdateUserByID.ExecContext")

		return err
	}

	return nil
}

func (repo *repository) GetUserByID(ctx context.Context, userID int64, lockOpt entities.LockingOpt) (userData entities.Users, err error) {
	if err := lockOpt.Validate(); err != nil {
		return userData, err
	}

	args := custom.Array{
		userID,
	}

	row := func(idx int) custom.Array {
		return custom.Array{
			&userData.UserID,
			&userData.Fullname,
			&userData.PhoneNumber,
			&userData.UserType,
			&userData.IsActive,
			&userData.CreatedAt,
		}

	}

	query := repository_query.GetUserByID

	if lockOpt.ForUpdate {
		query += " FOR UPDATE;"
	} else if lockOpt.ForUpdateNoWait {
		query += " FOR UPDATE NO WAIT;"
	}

	if err = new(database.SQL).Query(repo.conn.QueryContext(ctx, query, args...)).Scan(row); err != nil {
		repo.log.Z().Err(err).Msg("[repository]GetUserByID.QueryContext")
		return userData, err
	}

	return userData, err
}
