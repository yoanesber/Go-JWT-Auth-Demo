package entity

import (
	"time"

	"gopkg.in/go-playground/validator.v9"

	validation "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/validation-util"
)

// RefreshToken represents the refresh token entity in the database.
type RefreshToken struct {
	Token      string    `gorm:"column:token;type:text;primaryKey;unique;not null" json:"token" validate:"required"`
	UserID     int64     `gorm:"column:user_id;primaryKey;unique;not null" json:"userId" validate:"required"`
	User       *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user,omitempty"`
	ExpiryDate time.Time `gorm:"column:expiry_date;type:timestamptz;not null" json:"expiryDate" validate:"required"`
}

// RefreshTokenRequest represents the request payload for refreshing a token.
// It contains the refresh token that needs to be validated and used to obtain a new access token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// RefreshTokenResponse represents the response payload for refreshing a token.
// It contains the new access token, refresh token, expiration date, and token type.
type RefreshTokenResponse struct {
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	ExpirationDate string `json:"expirationDate"`
	TokenType      string `json:"tokenType"`
}

// TableName override the table name used by RefreshToken to `refresh_token`.
func (RefreshToken) TableName() string {
	return "refresh_token"
}

// Equals compares two RefreshToken objects for equality.
func (r *RefreshToken) Equals(other *RefreshToken) bool {
	if r == nil && other == nil {
		return true
	}

	if r == nil || other == nil {
		return false
	}

	if (r.Token != other.Token) ||
		(r.UserID != other.UserID) ||
		(r.ExpiryDate != other.ExpiryDate) {
		return false
	}

	return true
}

// Validate validates the RefreshTokenRequest struct using the validator package.
// It checks if the struct fields meet the specified validation rules.
func (a *RefreshTokenRequest) Validate() error {
	var v *validator.Validate = validation.GetValidator()

	if err := v.Struct(a); err != nil {
		return err
	}
	return nil
}
