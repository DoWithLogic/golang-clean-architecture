package usecase

import (
	"context"
	"database/sql"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users"
	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/errs"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability/instrumentation"
)

func (uc *usecase) UserUpdate(ctx context.Context, request dtos.UserUpdateRequest) error {
	ctx, span := instrumentation.NewTraceSpan(ctx, "UserUpdateUC")
	defer span.End()

	if _, err := uc.repo.UserDetail(ctx, entities.WithID(request.ID)); err != nil {
		return err
	}

	if request.ContactValue != nil {
		if uc.repo.IsUserExists(ctx, *request.ContactValue) {
			return errs.Conflict(errs.ErrUserAlreadyExists)
		}
	}

	var encryptedPassword *string
	if request.Password != nil {
		if newPassword := uc.crypto.EncodeSHA256(*request.Password); newPassword != "" {
			encryptedPassword = &newPassword
		}
	}

	err := uc.repo.WithTx(ctx, &sql.TxOptions{}, func(tx users.Repository) error {
		return tx.UpdateUser(ctx, request.ToUpdateUserEntity(encryptedPassword))
	})

	if err != nil {
		return errs.InternalServerError(err)
	}

	return nil
}
