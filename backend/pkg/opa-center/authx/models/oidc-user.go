package models

import "fmt"

type OIDCUser struct {
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	OriginalToken     string `json:"-"`
}

func (u *OIDCUser) GetAuthorizationHeader() string {
	return fmt.Sprintf("Bearer %s", u.OriginalToken)
}

func (u *OIDCUser) GetIdentifier() string {
	if u.PreferredUsername != "" {
		return u.PreferredUsername
	}

	return u.Email
}
