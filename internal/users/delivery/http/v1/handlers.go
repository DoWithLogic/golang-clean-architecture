package v1

import (
	"github.com/DoWithLogic/golang-clean-architecture/internal/users"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_jwt"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/apperror"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/observability/instrumentation"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/response"
	"github.com/labstack/echo/v4"
)

type handlers struct {
	uc users.Usecase
}

func NewHandlers(uc users.Usecase) *handlers {
	return &handlers{uc}
}

func (h *handlers) Login(c echo.Context) error {
	ctx, span := instrumentation.NewTraceSpan(c.Request().Context(), "LoginHandler")
	defer span.End()

	var request dtos.UserLoginRequest
	if err := c.Bind(&request); err != nil {
		return response.ErrorBuilder(apperror.BadRequest(err)).Send(c)
	}

	if err := request.Validate(); err != nil {
		return response.ErrorBuilder(apperror.BadRequest(err)).Send(c)
	}

	authData, err := h.uc.Login(ctx, request)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	return response.SuccessBuilder(authData).Send(c)
}

func (h *handlers) CreateUser(c echo.Context) error {
	ctx, span := instrumentation.NewTraceSpan(c.Request().Context(), "CreateUserHandler")
	defer span.End()

	var payload dtos.CreateUserRequest
	if err := c.Bind(&payload); err != nil {
		return response.ErrorBuilder(apperror.BadRequest(err)).Send(c)
	}

	if err := payload.Validate(); err != nil {
		return response.ErrorBuilder(apperror.BadRequest(err)).Send(c)
	}

	userID, err := h.uc.Create(ctx, payload)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	return response.SuccessBuilder(map[string]int64{"id": userID}).Send(c)
}

func (h *handlers) UserDetail(c echo.Context) error {
	ctx, span := instrumentation.NewTraceSpan(c.Request().Context(), "UserDetailHandler")
	defer span.End()

	userData, err := app_jwt.NewTokenInformation(c)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	data, err := h.uc.Detail(ctx, userData.Data.UserID)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	return response.SuccessBuilder(data).Send(c)
}

func (h *handlers) UpdateUser(c echo.Context) error {
	ctx, span := instrumentation.NewTraceSpan(c.Request().Context(), "UpdateUserHandler")
	defer span.End()

	var request dtos.UpdateUser
	if err := c.Bind(&request); err != nil {
		return response.ErrorBuilder(apperror.BadRequest(err)).Send(c)
	}

	if err := request.Validate(); err != nil {
		return response.ErrorBuilder(apperror.BadRequest(err)).Send(c)
	}

	userData, err := app_jwt.NewTokenInformation(c)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	args := dtos.UpdateUserRequest{
		UserID:     userData.Data.UserID,
		UpdateUser: request,
	}

	if err := h.uc.PartialUpdate(ctx, args); err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	return response.SuccessBuilder(nil).Send(c)
}

func (h *handlers) UpdateUserStatus(c echo.Context) error {
	ctx, span := instrumentation.NewTraceSpan(c.Request().Context(), "UpdateUserStatusHandler")
	defer span.End()

	var request dtos.UpdateUserStatus
	if err := c.Bind(&request); err != nil {
		return response.ErrorBuilder(apperror.BadRequest(err)).Send(c)
	}

	if err := request.Validate(); err != nil {
		return response.ErrorBuilder(apperror.BadRequest(err)).Send(c)
	}

	userData, err := app_jwt.NewTokenInformation(c)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	args := dtos.UpdateUserStatusRequest{
		UserID:           userData.Data.UserID,
		UpdateUserStatus: request,
	}

	if err := h.uc.UpdateStatus(ctx, args); err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	return response.SuccessBuilder(nil).Send(c)
}
