package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/internal/users/dtos"
	"github.com/DoWithLogic/golang-clean-architecture/internal/users/entities"
	usecases "github.com/DoWithLogic/golang-clean-architecture/internal/users/usecase"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
	"github.com/labstack/echo/v4"
)

type (
	Handlers interface {
		CreateUser(c echo.Context) error
		UpdateUser(c echo.Context) error
		UpdateUserStatus(c echo.Context) error
	}

	handlers struct {
		uc  usecases.Usecase
		log *zerolog.Logger
	}
)

const (
	BooleanTextTrue  = "true"
	BooleanTextFalse = "false"
)

var (
	ErrInvalidIsActive = errors.New("invalid is_active")
)

func NewHandlers(uc usecases.Usecase, log *zerolog.Logger) Handlers {
	return &handlers{uc, log}
}

func (h *handlers) CreateUser(c echo.Context) error {
	var (
		ctx, cancel = context.WithTimeout(c.Request().Context(), time.Duration(30*time.Second))
		payload     dtos.CreateUserPayload
	)
	defer cancel()

	if err := c.Bind(&payload); err != nil {
		h.log.Z().Err(err).Msg("[handlers]CreateUser.Bind")

		return c.JSON(http.StatusBadRequest, dtos.NewResponseError(
			http.StatusBadRequest,
			dtos.MsgFailed,
			err.Error()),
		)
	}

	if err := payload.Validate(); err != nil {
		h.log.Z().Err(err).Msg("[handlers]CreateUser.Validate")

		return c.JSON(http.StatusBadRequest, dtos.NewResponseError(
			http.StatusBadRequest,
			dtos.MsgFailed,
			err.Error()),
		)
	}

	argsCreateUser := entities.CreateUser{
		FUllName:    payload.FullName,
		PhoneNumber: payload.PhoneNumber,
	}

	createdID, err := h.uc.CreateUser(ctx, argsCreateUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dtos.NewResponseError(
			http.StatusInternalServerError,
			dtos.MsgFailed,
			err.Error()),
		)
	}

	return c.JSON(http.StatusOK, dtos.NewResponse(http.StatusOK, dtos.MsgSuccess, map[string]any{"user_id": createdID}))
}

func (h *handlers) UpdateUser(c echo.Context) error {
	var (
		ctx, cancel = context.WithTimeout(c.Request().Context(), time.Duration(30*time.Second))
		payload     dtos.UpdateUserPayload
	)
	defer cancel()

	h.log.Z().Info().Msg("[handlers]UpdateUser")

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Z().Err(err).Msg("[handlers]UpdateUser.ParseParam")

		return c.JSON(http.StatusBadRequest, dtos.NewResponseError(
			http.StatusBadRequest,
			dtos.MsgFailed,
			err.Error()),
		)
	}

	if err := c.Bind(&payload); err != nil {
		h.log.Z().Err(err).Msg("[handlers]UpdateUser.Bind")

		return c.JSON(http.StatusBadRequest, dtos.NewResponseError(
			http.StatusBadRequest,
			dtos.MsgFailed,
			err.Error()),
		)
	}

	if err := payload.Validate(); err != nil {
		h.log.Z().Err(err).Msg("[handlers]UpdateUser.Validate")

		return c.JSON(http.StatusBadRequest, dtos.NewResponseError(
			http.StatusBadRequest,
			dtos.MsgFailed,
			err.Error()),
		)
	}

	argsUpdateUser := entities.UpdateUsers{
		UserID:      userID,
		Fullname:    payload.Fullname,
		PhoneNumber: payload.PhoneNumber,
		UserType:    payload.UserType,
	}

	err = h.uc.UpdateUser(ctx, argsUpdateUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dtos.NewResponseError(
			http.StatusInternalServerError,
			dtos.MsgFailed,
			err.Error()),
		)
	}

	return c.JSON(http.StatusOK, dtos.NewResponse(http.StatusOK, dtos.MsgSuccess, nil))
}

func (h *handlers) UpdateUserStatus(c echo.Context) error {
	var (
		ctx, cancel = context.WithTimeout(c.Request().Context(), time.Duration(30*time.Second))
		payload     dtos.UpdateUserStatusPayload
	)
	defer cancel()

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.log.Z().Err(err).Msg("[handlers]UpdateUser.ParseParam")

		return c.JSON(http.StatusBadRequest, dtos.NewResponseError(
			http.StatusBadRequest,
			dtos.MsgFailed,
			err.Error()),
		)
	}

	switch c.QueryParam("is_active") {
	case BooleanTextFalse:
		payload.IsActive = false
	case BooleanTextTrue:
		payload.IsActive = true
	default:
		return c.JSON(http.StatusBadRequest, dtos.NewResponseError(
			http.StatusBadRequest,
			dtos.MsgFailed,
			ErrInvalidIsActive.Error()),
		)
	}

	argsUpdateUserStatus := entities.UpdateUserStatus{
		UserID:   userID,
		IsActive: payload.IsActive,
	}

	err = h.uc.UpdateUserStatus(ctx, argsUpdateUserStatus)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dtos.NewResponseError(
			http.StatusInternalServerError,
			dtos.MsgFailed,
			err.Error()),
		)
	}

	return c.JSON(http.StatusOK, dtos.NewResponse(http.StatusOK, dtos.MsgSuccess, nil))
}
