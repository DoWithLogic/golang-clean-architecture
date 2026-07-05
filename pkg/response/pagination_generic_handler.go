package response

import (
	"context"

	"github.com/labstack/echo/v4"
)

func GenericPaginationHandler[Request validatable, Response any](c echo.Context, request *Request, usecaseFn func(ctx context.Context, request *Request) (Response, error)) error {
	if err := c.Bind(request); err != nil {
		return ErrorBuilder(BadRequest(err)).Send(c)
	}

	if err := (*request).Validate(); err != nil {
		return ErrorBuilder(BadRequest(err)).Send(c)
	}

	result, err := usecaseFn(c.Request().Context(), request)
	if err != nil {
		return ErrorBuilder(err).Send(c)
	}

	return SuccessBuilder(result, *request).Send(c)
}
