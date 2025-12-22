package domain

import "context"

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, taskID string) (*Task, error)
	List(ctx context.Context) ([]Task, error)

	UpdateStatus(ctx context.Context, taskID string, status TaskStatus) error
}

type AssignmentRepository interface {
	Create(ctx context.Context, assignment *TaskAssignment) error

	GetByTask(ctx context.Context, taskID string) ([]TaskAssignment, error)

	GetByUser(ctx context.Context, userID string) ([]TaskAssignment, error)

	Delete(ctx context.Context, assignmentID string) error
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
