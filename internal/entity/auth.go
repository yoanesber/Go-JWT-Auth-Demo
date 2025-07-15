package entity

import (
	"gopkg.in/go-playground/validator.v9"

	validation "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/validation-util"
)

// LoginRequest represents the request payload for user login.
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

// LoginResponse represents the response payload for user login.
type LoginResponse struct {
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	ExpirationDate string `json:"expirationDate"`
	TokenType      string `json:"tokenType"`
}

// Validate validates the LoginRequest struct using the validator package.
// It checks if the struct fields meet the specified validation rules.
func (a *LoginRequest) Validate() error {
	var v *validator.Validate = validation.GetValidator()

	if err := v.Struct(a); err != nil {
		return err
	}
	return nil
}
