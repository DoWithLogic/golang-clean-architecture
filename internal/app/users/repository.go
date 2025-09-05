package users

import (
	"context"
	"database/sql"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/entities"
)

type Repository interface {
	WithTx(ctx context.Context, opt *sql.TxOptions, cb func(tx Repository) error) error

	AddUser(ctx context.Context, user *entities.User) error
	IsUserExists(ctx context.Context, contactValue string) bool
	UserDetail(ctx context.Context, opts ...entities.UserDetailOption) (user entities.User, err error)
	UpdateUser(ctx context.Context, user *entities.UpdateUser) error
}
