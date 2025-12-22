package domain

import "github.com/google/uuid"

type AssignmentScope struct {
	AssignmentID uuid.UUID
	TagID        uuid.UUID
}
