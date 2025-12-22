package domain

import (
	"time"

	"github.com/google/uuid"
)

type ActivityAction string

const (
	ActivityTaskCreated       ActivityAction = "task_created"
	ActivityAssignmentAdded   ActivityAction = "assignment_added"
	ActivityScopeAdded        ActivityAction = "scope_added"
	ActivityTaskStatusChanged ActivityAction = "task_status_changed"
)

type Activity struct {
	ID     uuid.UUID      `json:"id" db:"id"`
	TaskID uuid.UUID      `json:"task_id" db:"task_id"`
	UserID uuid.UUID      `json:"user_id" db:"user_id"`
	Action ActivityAction `json:"action" db:"action"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
