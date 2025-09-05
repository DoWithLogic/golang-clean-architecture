package usecase

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability/instrumentation"
)

func (uc *usecase) SignUp(ctx context.Context, request dtos.SignUpRequest) error {
	ctx, span := instrumentation.NewTraceSpan(ctx, "SignUpUC")
	defer span.End()

	if uc.repo.IsUserExists(ctx, request.ContactValue) {
		return apperror.Conflict(apperror.ErrUserAlreadyExists)
	}

	if err := uc.repo.AddUser(ctx, request.ToUserEntity(uc.crypto.EncodeSHA256(request.Password))); err != nil {
		return err
	}

	return nil
}
