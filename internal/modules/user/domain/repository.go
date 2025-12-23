package domain

import "github.com/google/uuid"

type UserRepository interface {
	GetAll() ([]User, error)

	GetByUsername(username string) (*User, error)

	GetByUserID(userID uuid.UUID) (*User, error)

	Create(user *User) error

	Delete(id uuid.UUID) error
}
