package usecase

import (
	"context"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability/instrumentation"
)

func (uc *usecase) Login(ctx context.Context, request dtos.UserLoginRequest) (response dtos.UserLoginResponse, err error) {
	ctx, span := instrumentation.NewTraceSpan(ctx, "LoginUC")
	defer span.End()

	userData, err := uc.repo.UserDetail(ctx, entities.WithContactValue(request.ContactValue))
	if err != nil {
		return response, err
	}

	if !userData.IsPasswordValid(uc.crypto.EncodeSHA256(request.Password)) {
		return response, apperror.Unauthorized(apperror.ErrInvalidPassword)
	}

	expiredAt := time.Now().Add(time.Minute * 60).Unix()
	jwtToken, err := uc.appJwt.GenerateToken(ctx, userData.ToJWTData(expiredAt))
	if err != nil {
		return response, apperror.InternalServerError(err)
	}

	return dtos.ToUserLoginResponse(jwtToken, expiredAt-time.Now().Unix()), nil
}
