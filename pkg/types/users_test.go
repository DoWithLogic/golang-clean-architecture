package types

import "testing"

func TestContactTypeConstants(t *testing.T) {
	tests := []struct {
		name string
		got  CONTACT_TYPE
		want string
	}{
		{
			name: "email",
			got:  CONTACT_TYPE_EMAIL,
			want: "EMAIL",
		},
		{
			name: "phone",
			got:  CONTACT_TYPE_PHONE,
			want: "PHONE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.want {
				t.Errorf("expected %q, got %q", tt.want, tt.got)
			}
		})
	}
}

func TestUserStatusConstants(t *testing.T) {
	tests := []struct {
		name string
		got  USER_STATUS
		want string
	}{
		{
			name: "pending",
			got:  PENDING,
			want: "PENDING",
		},
		{
			name: "active",
			got:  ACTIVE,
			want: "ACTIVE",
		},
		{
			name: "reject",
			got:  REJECT,
			want: "REJECT",
		},
		{
			name: "banned",
			got:  BANNED,
			want: "CLOSED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.want {
				t.Errorf("expected %q, got %q", tt.want, tt.got)
			}
		})
	}
}

func TestLanguageConstants(t *testing.T) {
	tests := []struct {
		name string
		got  LANGUAGE
		want string
	}{
		{
			name: "english",
			got:  LANGUAGE_EN,
			want: "EN",
		},
		{
			name: "indonesian",
			got:  LANGUAGE_ID,
			want: "ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.want {
				t.Errorf("expected %q, got %q", tt.want, tt.got)
			}
		})
	}
}
