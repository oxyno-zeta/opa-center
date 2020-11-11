package models

import (
	"testing"
)

func TestOIDCUser_GetIdentifier(t *testing.T) {
	type fields struct {
		PreferredUsername string
		Name              string
		GivenName         string
		FamilyName        string
		Email             string
		EmailVerified     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "all empty",
			fields: fields{
				PreferredUsername: "",
				Email:             "",
			},
			want: "",
		},
		{
			name: "empty email",
			fields: fields{
				PreferredUsername: "username",
				Email:             "",
			},
			want: "username",
		},
		{
			name: "empty username",
			fields: fields{
				PreferredUsername: "",
				Email:             "email",
			},
			want: "email",
		},
		{
			name: "all set",
			fields: fields{
				PreferredUsername: "username",
				Email:             "email",
			},
			want: "username",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &OIDCUser{
				PreferredUsername: tt.fields.PreferredUsername,
				Name:              tt.fields.Name,
				GivenName:         tt.fields.GivenName,
				FamilyName:        tt.fields.FamilyName,
				Email:             tt.fields.Email,
				EmailVerified:     tt.fields.EmailVerified,
			}
			if got := u.GetIdentifier(); got != tt.want {
				t.Errorf("OIDCUser.GetIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOIDCUser_GetAuthorizationHeader(t *testing.T) {
	type fields struct {
		OriginalToken string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Empty case",
			fields: fields{OriginalToken: ""},
			want:   "Bearer ",
		},
		{
			name:   "Normal case",
			fields: fields{OriginalToken: "fake"},
			want:   "Bearer fake",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &OIDCUser{
				OriginalToken: tt.fields.OriginalToken,
			}
			if got := u.GetAuthorizationHeader(); got != tt.want {
				t.Errorf("OIDCUser.GetAuthorizationHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
