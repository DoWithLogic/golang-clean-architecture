package response

// ResponseMessage defines standard response messages.
type ResponseMessage string

const (
	InternalServerErrorMessage ResponseMessage = "internal_server_error"
	BadRequestMessage          ResponseMessage = "bad_request"
	SuccessMessage             ResponseMessage = "success"
	UnauthorizedMessage        ResponseMessage = "unauthorized"
	ForbiddenMessage           ResponseMessage = "forbidden"
	NotFoundMessage            ResponseMessage = "not_found"
	ConflictMessage            ResponseMessage = "conflict"
	GatewayTimeOutMessage      ResponseMessage = "gateway_timeout"
	TooManyRequestsMessage     ResponseMessage = "too_many_requests"
)

func (rm ResponseMessage) String() string {
	return string(rm)
}
