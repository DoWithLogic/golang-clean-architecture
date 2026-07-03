package usecase

import (
	"context"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/internal/app/users/entities"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability/instrumentation"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response/app_error"
)

func (uc *usecase) Login(ctx context.Context, request dtos.UserLoginRequest) (result dtos.UserLoginResponse, err error) {
	ctx, span := instrumentation.NewTraceSpan(ctx, "LoginUC")
	defer span.End()

	userData, err := uc.repo.UserDetail(ctx, entities.WithContactValue(request.ContactValue))
	if err != nil {
		return result, err
	}

	if !userData.IsPasswordValid(uc.crypto.EncodeSHA256(request.Password)) {
		return result, response.Unauthorized(app_error.ErrInvalidPassword)
	}

	expiredAt := time.Now().Add(time.Minute * 60).Unix()
	jwtToken, err := uc.appJwt.CreateJWT(userData.ToJWTData(expiredAt))
	if err != nil {
		return result, response.InternalServerError(err)
	}

	return dtos.ToUserLoginResponse(jwtToken, expiredAt-time.Now().Unix()), nil
}
