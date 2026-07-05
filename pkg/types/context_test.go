package types

import "testing"

func TestContextKey_String(t *testing.T) {
	tests := []struct {
		name string
		key  CONTEXT_KEY
		want string
	}{
		{
			name: "credential data context key",
			key:  CredentialDataContextKey,
			want: "credential_data",
		},
		{
			name: "custom context key",
			key:  CONTEXT_KEY("custom_key"),
			want: "custom_key",
		},
		{
			name: "empty context key",
			key:  CONTEXT_KEY(""),
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
