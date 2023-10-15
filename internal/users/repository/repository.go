package repository

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/repository/repository_query"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasource"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/utils"
)

type (
	Repository interface {
		SaveNewUser(context.Context, entities.Users) (int64, error)
		UpdateUserByID(context.Context, entities.UpdateUsers) error
		GetUserByID(context.Context, int64, ...entities.LockingOpt) (entities.Users, error)
		UpdateUserStatusByID(context.Context, entities.UpdateUserStatus) error
	}

	repository struct {
		conn datasource.SQLTxConn
		log  *zerolog.Logger
	}
)

func NewRepository(conn datasource.SQLTxConn, log *zerolog.Logger) Repository {
	return &repository{conn, log}
}

func (repo *repository) SaveNewUser(ctx context.Context, user entities.Users) (int64, error) {
	args := utils.Array{
		user.Fullname,
		user.PhoneNumber,
		user.IsActive,
		user.UserType,
		user.CreatedAt,
		user.CreatedBy,
	}

	var userID int64
	err := new(datasource.SQL).Exec(repo.conn.ExecContext(ctx, repository_query.InsertUsers, args...)).Scan(nil, &userID)
	if err != nil {
		repo.log.Z().Err(err).Msg("[repository]SaveNewUser.ExecContext")

		return userID, err
	}

	return userID, err
}

func (repo *repository) UpdateUserByID(ctx context.Context, user entities.UpdateUsers) error {
	args := utils.Array{
		user.Fullname, user.Fullname,
		user.PhoneNumber, user.PhoneNumber,
		user.UserType, user.UserType,
		user.UpdatedAt,
		user.UpdatedBy,
		user.UserID,
	}

	err := new(datasource.SQL).Exec(repo.conn.ExecContext(ctx, repository_query.UpdateUsers, args...)).Scan(nil, nil)
	if err != nil {
		repo.log.Z().Err(err).Msg("[repository]UpdateUserByID.ExecContext")

		return err
	}

	return nil
}

func (repo *repository) GetUserByID(ctx context.Context, userID int64, lockOpt ...entities.LockingOpt) (userData entities.Users, err error) {
	args := utils.Array{
		userID,
	}

	row := func(idx int) utils.Array {
		return utils.Array{
			&userData.UserID,
			&userData.Fullname,
			&userData.PhoneNumber,
			&userData.UserType,
			&userData.IsActive,
			&userData.CreatedAt,
		}
	}

	query := repository_query.GetUserByID

	if len(lockOpt) >= 1 {
		if lockOpt[0].ForUpdate {
			query += " FOR UPDATE;"
		} else {
			query += " FOR UPDATE NO WAIT;"
		}
	}

	if err = new(datasource.SQL).Query(repo.conn.QueryContext(ctx, query, args...)).Scan(row); err != nil {
		repo.log.Z().Err(err).Msg("[repository]GetUserByID.QueryContext")
		return userData, err
	}

	return userData, err
}

func (repo *repository) UpdateUserStatusByID(ctx context.Context, req entities.UpdateUserStatus) error {
	args := utils.Array{
		req.IsActive,
		req.UpdatedAt,
		req.UpdatedBy,
		req.UserID,
	}

	var updatedID int64
	err := new(datasource.SQL).Exec(repo.conn.ExecContext(ctx, repository_query.UpdateUserStatusByID, args...)).Scan(nil, &updatedID)
	if err != nil {
		repo.log.Z().Err(err).Msg("[repository]UpdateUserStatusByID.ExecContext")

		return err
	}

	return nil
}
