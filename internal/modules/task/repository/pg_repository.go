package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresTaskRepository struct {
	db *sqlx.DB
}

func NewPostgresTaskRepository(db *sqlx.DB) domain.TaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, title, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		task.ID, task.Title, task.Status, task.CreatedBy, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, taskID string) (*domain.Task, error) {
	task := &domain.Task{}
	query := `SELECT id, title, status, created_by, created_at, updated_at FROM tasks WHERE id = $1`
	err := r.db.GetContext(ctx, task, query, taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return task, nil
}

func (r *PostgresTaskRepository) UpdateStatus(ctx context.Context, taskID string, status domain.TaskStatus) error {
	query := `UPDATE tasks SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), taskID)
	return err
}

func (r *PostgresTaskRepository) List(ctx context.Context) ([]domain.Task, error) {
	tasks := []domain.Task{}
	query := `SELECT id, title, status, created_by, created_at, updated_at FROM tasks ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &tasks, query)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

type PostgresAssignmentRepository struct {
	db *sqlx.DB
}

func NewPostgresAssignmentRepository(db *sqlx.DB) domain.AssignmentRepository {
	return &PostgresAssignmentRepository{db: db}
}

func (r *PostgresAssignmentRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r *PostgresAssignmentRepository) Create(ctx context.Context, tx *sqlx.Tx, assignment *domain.TaskAssignment) error {
	query := `
		INSERT INTO task_assignments (id, task_id, user_id, created_at)
		VALUES ($1, $2, $3, $4)
	`

	var executor sqlx.ExtContext = r.db
	if tx != nil {
		executor = tx
	}

	_, err := executor.ExecContext(ctx, query,
		assignment.ID, assignment.TaskID, assignment.UserID, assignment.CreatedAt)
	return err
}

func (r *PostgresAssignmentRepository) GetByTask(ctx context.Context, taskID string) ([]domain.TaskAssignment, error) {
	assignments := []domain.TaskAssignment{}
	query := `SELECT id, task_id, user_id, created_at FROM task_assignments WHERE task_id = $1`
	err := r.db.SelectContext(ctx, &assignments, query, taskID)
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

func (r *PostgresAssignmentRepository) GetByUser(ctx context.Context, userID string) ([]domain.TaskAssignment, error) {
	assignments := []domain.TaskAssignment{}
	query := `SELECT id, task_id, user_id, created_at FROM task_assignments WHERE user_id = $1`
	err := r.db.SelectContext(ctx, &assignments, query, userID)
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

func (r *PostgresAssignmentRepository) Delete(ctx context.Context, assignmentID string) error {
	query := `DELETE FROM task_assignments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, assignmentID)
	return err
}

type PostgresScopeRepository struct {
	db *sqlx.DB
}

func NewPostgresScopeRepository(db *sqlx.DB) domain.ScopeRepository {
	return &PostgresScopeRepository{db: db}
}

func (r *PostgresScopeRepository) AddToAssignment(ctx context.Context, assignmentID string, scopeID string) error {
	query := `INSERT INTO assignment_scopes (assignment_id, scope_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, assignmentID, scopeID)
	return err
}

func (r *PostgresScopeRepository) RemoveFromAssignment(ctx context.Context, assignmentID string, scopeID string) error {
	query := `DELETE FROM assignment_scopes WHERE assignment_id = $1 AND scope_id = $2`
	_, err := r.db.ExecContext(ctx, query, assignmentID, scopeID)
	return err
}

func (r *PostgresScopeRepository) GetByAssignment(ctx context.Context, assignmentID string) ([]domain.Scope, error) {
	scopes := []domain.Scope{}
	query := `
		SELECT s.id, s.name
		FROM scopes s
		INNER JOIN assignment_scopes as ON as.scope_id = s.id
		WHERE as.assignment_id = $1
	`
	err := r.db.SelectContext(ctx, &scopes, query, assignmentID)
	if err != nil {
		return nil, err
	}
	return scopes, nil
}

type PostgresScopeLookupRepository struct {
	db *sqlx.DB
}

func NewPostgresScopeLookupRepository(db *sqlx.DB) domain.ScopeLookupRepository {
	return &PostgresScopeLookupRepository{db: db}
}

func (r *PostgresScopeLookupRepository) GetAll(ctx context.Context) ([]domain.Scope, error) {
	scopes := []domain.Scope{}
	query := `SELECT id, name FROM scopes ORDER BY name ASC`
	err := r.db.SelectContext(ctx, &scopes, query)
	if err != nil {
		return nil, err
	}
	return scopes, nil
}

func (r *PostgresScopeLookupRepository) GetByID(ctx context.Context, scopeID string) (*domain.Scope, error) {
	scope := &domain.Scope{}
	query := `SELECT id, name FROM scopes WHERE id = $1`
	err := r.db.GetContext(ctx, scope, query, scopeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return scope, nil
}

type PostgresActivityRepository struct {
	db *sqlx.DB
}

func NewPostgresActivityRepository(db *sqlx.DB) domain.ActivityRepository {
	return &PostgresActivityRepository{db: db}
}

func (r *PostgresActivityRepository) Create(ctx context.Context, activity *domain.Activity) error {
	query := `
		INSERT INTO task_activities (id, task_id, user_id, action, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		activity.ID, activity.TaskID, activity.UserID, activity.Action, activity.CreatedAt)
	return err
}

func (r *PostgresActivityRepository) GetByTask(ctx context.Context, taskID string) ([]domain.Activity, error) {
	activities := []domain.Activity{}
	query := `SELECT id, task_id, user_id, action, created_at FROM task_activities WHERE task_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &activities, query, taskID)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (r *PostgresActivityRepository) GetByUser(ctx context.Context, userID string) ([]domain.Activity, error) {
	activities := []domain.Activity{}
	query := `SELECT id, task_id, user_id, action, created_at FROM task_activities WHERE user_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &activities, query, userID)
	if err != nil {
		return nil, err
	}
	return activities, nil
}
