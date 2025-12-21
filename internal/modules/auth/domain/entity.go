package domain

import "github.com/golang-jwt/jwt/v5"

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Ad       string `json:"ad"`
	Soyad    string `json:"soyad"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type UserCredentials struct {
	Username string
	Password string
	Role     string
	Ad       string
	Soyad    string
}

type ErrInvalidCredentials struct{}

func (e ErrInvalidCredentials) Error() string {
	return "invalid username or password"
}

type ErrTokenGeneration struct{}

func (e ErrTokenGeneration) Error() string {
	return "failed to generate token"
}
