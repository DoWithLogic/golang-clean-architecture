package types

import "testing"

func TestHeaderKey_String(t *testing.T) {
	tests := []struct {
		name string
		key  HEADER_KEY
		want string
	}{
		{
			name: "authorization header",
			key:  AuthorizationHeaderKey,
			want: "Authorization",
		},
		{
			name: "custom header",
			key:  HEADER_KEY("X-Custom-Header"),
			want: "X-Custom-Header",
		},
		{
			name: "empty header",
			key:  HEADER_KEY(""),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.key.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
