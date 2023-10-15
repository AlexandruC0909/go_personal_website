package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
