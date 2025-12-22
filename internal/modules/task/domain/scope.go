package domain

import "github.com/google/uuid"

type Scope struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}
