package usecase

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability/instrumentation"
)

func (uc *usecase) TransitionUserStatus(ctx context.Context, request dtos.TransitionUserStatusRequest) error {
	ctx, span := instrumentation.NewTraceSpan(ctx, "TransitionUserStatusUC")
	defer span.End()

	if _, err := uc.repo.UserDetail(ctx, entities.WithID(request.ID)); err != nil {
		return err
	}

	return uc.repo.UpdateUser(ctx, request.ToUpdateStatusEntity())
}
