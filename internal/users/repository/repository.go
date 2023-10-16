package repository

import (
	"context"
	"database/sql"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/repository/repository_query"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasource"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type (
	Repository interface {
		Atomic(ctx context.Context, opt *sql.TxOptions, repo func(tx Repository) error) error

		SaveNewUser(context.Context, entities.Users) (int64, error)
		UpdateUserByID(context.Context, entities.UpdateUsers) error
		GetUserByID(context.Context, int64, ...entities.LockingOpt) (entities.Users, error)
		UpdateUserStatusByID(context.Context, entities.UpdateUserStatus) error
	}

	repository struct {
		db   *sqlx.DB
		conn datasource.ConnTx
		log  *zerolog.Logger
	}
)

func NewRepository(c *sqlx.DB, l *zerolog.Logger) Repository {
	return &repository{conn: c, log: l, db: c}
}

// Atomic implements vendor.Repository for transaction query
func (r *repository) Atomic(ctx context.Context, opt *sql.TxOptions, repo func(tx Repository) error) error {
	txConn, err := r.db.BeginTx(ctx, opt)
	if err != nil {
		r.log.Z().Err(err).Msg("[repository]Atomic.BeginTxx")

		return err
	}

	newRepository := &repository{conn: txConn, db: r.db}

	repo(newRepository)

	if err := new(datasource.DataSource).EndTx(txConn, err); err != nil {
		return err
	}

	return nil
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
	err := new(datasource.DataSource).ExecSQL(repo.conn.ExecContext(ctx, repository_query.InsertUsers, args...)).Scan(nil, &userID)
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

	err := new(datasource.DataSource).ExecSQL(repo.conn.ExecContext(ctx, repository_query.UpdateUsers, args...)).Scan(nil, nil)
	if err != nil {
		repo.log.Z().Err(err).Msg("[repository]UpdateUserByID.ExecContext")

		return err
	}

	return nil
}

func (repo *repository) GetUserByID(ctx context.Context, userID int64, options ...entities.LockingOpt) (userData entities.Users, err error) {
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

	if len(options) >= 1 && options[0].PessimisticLocking {
		query += " FOR UPDATE"
	}

	if err = new(datasource.DataSource).QuerySQL(repo.conn.QueryContext(ctx, query, args...)).Scan(row); err != nil {
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
	err := new(datasource.DataSource).ExecSQL(repo.conn.ExecContext(ctx, repository_query.UpdateUserStatusByID, args...)).Scan(nil, &updatedID)
	if err != nil {
		repo.log.Z().Err(err).Msg("[repository]UpdateUserStatusByID.ExecContext")

		return err
	}

	return nil
}
