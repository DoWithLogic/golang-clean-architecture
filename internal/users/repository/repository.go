package repository

import (
	"context"
	"database/sql"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/repository/repository_query"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/datasource"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db   *sqlx.DB
	conn datasource.ConnTx
}

func NewRepository(c *sqlx.DB) users.Repository {
	return &repository{conn: c, db: c}
}

// Atomic implements Repository Interface for transaction query
func (r *repository) Atomic(ctx context.Context, opt *sql.TxOptions, repo func(tx users.Repository) error) error {
	txConn, err := r.db.BeginTxx(ctx, opt)
	if err != nil {
		return err
	}

	newRepository := &repository{conn: txConn, db: r.db}

	repo(newRepository)

	if err := new(datasource.DataSource).EndTx(txConn, err); err != nil {
		return err
	}

	return nil
}

func (repo *repository) SaveNewUser(ctx context.Context, user entities.Users) (userID int64, err error) {
	args := utils.Array{
		user.Email,
		user.Password,
		user.Fullname,
		user.PhoneNumber,
		user.UserType,
		user.IsActive,
		user.CreatedAt,
		user.CreatedBy,
	}

	err = new(datasource.DataSource).ExecSQL(repo.conn.ExecContext(ctx, repository_query.InsertUsers, args...)).Scan(nil, &userID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}

func (repo *repository) UpdateUserByID(ctx context.Context, user entities.UpdateUsers) error {
	args := utils.Array{
		user.Email, user.Email,
		user.Fullname, user.Fullname,
		user.PhoneNumber, user.PhoneNumber,
		user.UserType, user.UserType,
		user.UpdatedAt,
		user.UpdatedBy,
		user.UserID,
	}

	err := new(datasource.DataSource).ExecSQL(repo.conn.ExecContext(ctx, repository_query.UpdateUsers, args...)).Scan(nil, nil)
	if err != nil {
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
			&userData.Email,
			&userData.Fullname,
			&userData.PhoneNumber,
			&userData.UserType,
			&userData.IsActive,
			&userData.CreatedAt,
			&userData.CreatedBy,
		}
	}

	query := repository_query.GetUserByID

	if len(options) >= 1 && options[0].PessimisticLocking {
		query += " FOR UPDATE"
	}

	if err = new(datasource.DataSource).QuerySQL(repo.conn.QueryxContext(ctx, query, args...)).Scan(row); err != nil {
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
		return err
	}

	return nil
}

func (repo *repository) IsUserExist(ctx context.Context, email string) bool {
	args := utils.Array{email}

	var id int64
	row := func(idx int) utils.Array {
		return utils.Array{
			&id,
		}
	}

	err := new(datasource.DataSource).QuerySQL(repo.conn.QueryxContext(ctx, repository_query.IsUserExist, args...)).Scan(row)
	if err != nil {
		return false
	}

	return id != 0
}

func (repo *repository) GetUserByEmail(ctx context.Context, email string) (userDetail entities.Users, err error) {
	args := utils.Array{
		email,
	}

	row := func(idx int) utils.Array {
		return utils.Array{
			&userDetail.UserID,
			&userDetail.Email,
			&userDetail.Password,
		}
	}

	err = new(datasource.DataSource).QuerySQL(repo.conn.QueryxContext(ctx, repository_query.GetUserByEmail, args...)).Scan(row)
	if err != nil {
		return entities.Users{}, err
	}

	return userDetail, err
}
