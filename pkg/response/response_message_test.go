package response

import "testing"

func TestResponseMessage_String(t *testing.T) {
	tests := []struct {
		name string
		msg  ResponseMessage
		want string
	}{
		{
			name: "internal server error",
			msg:  InternalServerErrorMessage,
			want: "internal_server_error",
		},
		{
			name: "bad request",
			msg:  BadRequestMessage,
			want: "bad_request",
		},
		{
			name: "success",
			msg:  SuccessMessage,
			want: "success",
		},
		{
			name: "unauthorized",
			msg:  UnauthorizedMessage,
			want: "unauthorized",
		},
		{
			name: "forbidden",
			msg:  ForbiddenMessage,
			want: "forbidden",
		},
		{
			name: "not found",
			msg:  NotFoundMessage,
			want: "not_found",
		},
		{
			name: "conflict",
			msg:  ConflictMessage,
			want: "conflict",
		},
		{
			name: "gateway timeout",
			msg:  GatewayTimeOutMessage,
			want: "gateway_timeout",
		},
		{
			name: "too many requests",
			msg:  TooManyRequestsMessage,
			want: "too_many_requests",
		},
		{
			name: "custom message",
			msg:  ResponseMessage("custom"),
			want: "custom",
		},
		{
			name: "empty message",
			msg:  ResponseMessage(""),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
