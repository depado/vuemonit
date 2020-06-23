package router

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type RefreshQuery struct {
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshQuery) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RefreshToken, validation.Required),
	)
}

type CredentialQuery struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c CredentialQuery) ValidateRegister() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required, validation.Length(12, 50)),
	)
}

func (c CredentialQuery) ValidateLogin() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required),
		validation.Field(&c.Password, validation.Required),
	)
}
