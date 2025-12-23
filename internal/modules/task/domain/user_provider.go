package domain

import (
	"github.com/google/uuid"
)

type UserProvider interface {
	GetUserByID(userID uuid.UUID) (*UserInfo, error)
}

type UserInfo struct {
	ID       uuid.UUID
	Username string
	Email    string
}
