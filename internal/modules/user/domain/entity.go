package domain

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username" validate:"required"`
	Password string    `json:"password,omitempty" db:"password" validate:"required"`
	Role     string    `json:"role" db:"role" validate:"required"`
	Ad       string    `json:"ad" db:"ad"`
	Soyad    string    `json:"soyad" db:"soyad"`
	Telefon  string    `json:"telefon" db:"telefon"`
	Email    string    `json:"email" db:"email"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=3"`
	Role     string `json:"role" validate:"required"`
	Ad       string `json:"ad"`
	Soyad    string `json:"soyad"`
	Telefon  string `json:"telefon"`
	Email    string `json:"email"`
}

type ErrUserNotFound struct{}

func (e ErrUserNotFound) Error() string {
	return "user not found"
}

type ErrUsernameTaken struct{}

func (e ErrUsernameTaken) Error() string {
	return "username already taken"
}
