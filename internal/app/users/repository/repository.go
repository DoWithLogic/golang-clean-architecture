package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users"
	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability/instrumentation"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) WithTx(ctx context.Context, opt *sql.TxOptions, cb func(tx users.Repository) error) error {
	tx := r.db.Begin(opt)

	if err := cb(&repository{db: tx}); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (r *repository) AddUser(ctx context.Context, user *entities.User) error {
	ctx, span := instrumentation.NewTraceSpan(ctx, "AddUserRepo")
	defer span.End()

	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) IsUserExists(ctx context.Context, contactValue string) bool {
	ctx, span := instrumentation.NewTraceSpan(ctx, "IsUserExistsRepo")
	defer span.End()

	var count int64
	if err := r.db.WithContext(ctx).Table(types.TABLE_NAME_USERS.String()).Where("contact_value = ?", contactValue).Count(&count).Error; err != nil {
		return false
	}

	return count > 0
}

func (r *repository) UserDetail(ctx context.Context, opts ...entities.UserDetailOption) (user entities.User, err error) {
	ctx, span := instrumentation.NewTraceSpan(ctx, "UserDetailRepo")
	defer span.End()

	request := new(entities.UserDetailRequest)
	for _, opt := range opts {
		opt.Apply(request)
	}

	baseQuery := r.db.WithContext(ctx).Table(types.TABLE_NAME_USERS.String())
	if request.ID != nil {
		baseQuery = baseQuery.Where("id = ?", request.ID)
	}

	if request.ContactValue != nil {
		baseQuery = baseQuery.Where("contact_value = ?", request.ContactValue)
	}

	if err := baseQuery.Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperror.NotFound(apperror.ErrUserNotFound)
		}

		return user, err
	}

	return user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *entities.UpdateUser) error {
	ctx, span := instrumentation.NewTraceSpan(ctx, "UpdateUserRepo")
	defer span.End()

	return r.db.WithContext(ctx).Model(&entities.User{}).Where("id = ?", user.ID).Updates(user).Error
}
