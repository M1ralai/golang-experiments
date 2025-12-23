package domain

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, taskID string) (*Task, error)
	List(ctx context.Context) ([]Task, error)

	UpdateStatus(ctx context.Context, taskID string, status TaskStatus) error
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
}

type AssignmentRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, assignment *TaskAssignment) error
	Delete(ctx context.Context, assignmentID string) error
	GetByTask(ctx context.Context, taskID string) ([]TaskAssignment, error)
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
}

type ScopeRepository interface {
	AddToAssignment(ctx context.Context, assignmentID string, scopeID string) error

	RemoveFromAssignment(ctx context.Context, assignmentID string, scopeID string) error

	GetByAssignment(ctx context.Context, assignmentID string) ([]Scope, error)
}

type ScopeLookupRepository interface {
	GetAll(ctx context.Context) ([]Scope, error)
	GetByID(ctx context.Context, scopeID string) (*Scope, error)
}

type ActivityRepository interface {
	Create(ctx context.Context, activity *Activity) error

	GetByTask(ctx context.Context, taskID string) ([]Activity, error)

	GetByUser(ctx context.Context, userID string) ([]Activity, error)
}
