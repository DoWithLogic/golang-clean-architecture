package users

import (
	"context"
	"database/sql"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
)

type Repository interface {
	Atomic(ctx context.Context, opt *sql.TxOptions, repo func(tx Repository) error) error

	GetUserByID(context.Context, int64, ...entities.LockingOpt) (entities.User, error)
	GetUserByEmail(context.Context, string) (entities.User, error)
	SaveNewUser(context.Context, entities.User) (int64, error)
	UpdateUserByID(context.Context, entities.UpdateUser) error
	UpdateUserStatusByID(context.Context, entities.UpdateUserStatus) error
	IsUserExist(ctx context.Context, email string) bool
}
