package usecase

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability/instrumentation"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response/app_error"
)

func (uc *usecase) SignUp(ctx context.Context, request dtos.SignUpRequest) error {
	ctx, span := instrumentation.NewTraceSpan(ctx, "SignUpUC")
	defer span.End()

	if uc.repo.IsUserExists(ctx, request.ContactValue) {
		return response.Conflict(app_error.ErrUserAlreadyExists)
	}

	if err := uc.repo.AddUser(ctx, request.ToUserEntity(uc.crypto.EncodeSHA256(request.Password))); err != nil {
		return err
	}

	return nil
}
