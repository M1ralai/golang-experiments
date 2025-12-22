package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

type Task struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	Status    TaskStatus `json:"status" db:"status"`
	CreatedBy uuid.UUID  `json:"created_by" db:"created_by"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Request DTOs
type CreateTaskRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

type UpdateStatusRequest struct {
	Status TaskStatus `json:"status" validate:"required,oneof=todo in_progress done"`
}

type AssignTaskRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}
