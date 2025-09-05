package usecase

import (
	"context"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability/instrumentation"
)

func (uc *usecase) UserDetail(ctx context.Context, request dtos.UserDetailRequest) (userData dtos.User, err error) {
	ctx, span := instrumentation.NewTraceSpan(ctx, "UserDetailUC")
	defer span.End()

	userDetail, err := uc.repo.UserDetail(ctx, request.ToUserDetailOption())
	if err != nil {
		return userData, err
	}

	return dtos.ToUserDTO(userDetail), nil
}
