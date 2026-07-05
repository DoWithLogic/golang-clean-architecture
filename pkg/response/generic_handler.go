package response

import (
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type validatable interface {
	Validate() error
}

func GenericHandler[Request validatable](c echo.Context, request *Request, usecaseFn any) error {
	if err := c.Bind(request); err != nil {
		return ErrorBuilder(BadRequest(err)).Send(c)
	}

	if err := (*request).Validate(); err != nil {
		return ErrorBuilder(BadRequest(err)).Send(c)
	}

	fnValue := reflect.ValueOf(usecaseFn)
	fnType := fnValue.Type()

	if fnType.Kind() != reflect.Func {
		return ErrorBuilder(InternalServerError(errors.New("usecaseFn is not a function"))).Send(c)
	}

	// check input param
	if fnType.NumIn() != 2 {
		return ErrorBuilder(InternalServerError(errors.New("function must have 2 input parameters"))).Send(c)
	}

	// Check output number
	if fnType.NumOut() < 1 || fnType.NumOut() > 2 {
		return ErrorBuilder(InternalServerError(errors.New("function must have 1 or 2 return values"))).Send(c)
	}

	// Prepare arguments for the function call
	args := []reflect.Value{
		reflect.ValueOf(c.Request().Context()),
		reflect.ValueOf(*request),
	}

	// Call the function
	results := fnValue.Call(args)

	// Handle return values based on the number of results
	switch len(results) {
	case 1: // HandlerFn: returns only an error
		err := results[0].Interface()
		if err != nil {
			return ErrorBuilder(err.(error)).Send(c)
		}

		return SuccessBuilder(nil).Send(c)

	case 2: // HandlerWithResponseFn: returns a response and an error
		response := results[0].Interface()
		err := results[1].Interface()

		if err != nil {
			return ErrorBuilder(err.(error)).Send(c)
		}

		return SuccessBuilder(response).Send(c)

	default:
		return ErrorBuilder(InternalServerError(errors.New("unsupported function signature"))).Send(c)
	}
}
