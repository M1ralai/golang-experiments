package outbox

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OutboxEvent struct {
	ID            uuid.UUID       `db:"id"`
	AggregateType string          `db:"aggregate_type"`
	AggregateID   uuid.UUID       `db:"aggregate_id"`
	EventType     string          `db:"event_type"`
	Payload       json.RawMessage `db:"payload"`
	CreatedAt     time.Time       `db:"created_at"`
	ProcessedAt   *time.Time      `db:"processed_at"`
	RetryCount    int             `db:"retry_count"`
	LastError     *string         `db:"last_error"`
}

type Repository interface {
	Create(ctx context.Context, tx *sqlx.Tx, event *OutboxEvent) error
	GetUnprocessed(ctx context.Context, limit int) ([]*OutboxEvent, error)
	MarkProcessed(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, errorMsg string) error
}

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, tx *sqlx.Tx, event *OutboxEvent) error {
	query := `
		INSERT INTO outbox_events (aggregate_type, aggregate_id, event_type, payload)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	var executor sqlx.ExtContext = r.db
	if tx != nil {
		executor = tx
	}

	return executor.QueryRowxContext(
		ctx, query,
		event.AggregateType,
		event.AggregateID,
		event.EventType,
		event.Payload,
	).Scan(&event.ID, &event.CreatedAt)
}

func (r *postgresRepository) GetUnprocessed(ctx context.Context, limit int) ([]*OutboxEvent, error) {
	query := `
		SELECT id, aggregate_type, aggregate_id, event_type, payload,
		       created_at, processed_at, retry_count, last_error
		FROM outbox_events
		WHERE processed_at IS NULL AND retry_count < 5
		ORDER BY created_at ASC
		LIMIT $1
	`

	var events []*OutboxEvent
	err := r.db.SelectContext(ctx, &events, query, limit)
	return events, err
}

func (r *postgresRepository) MarkProcessed(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE outbox_events SET processed_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresRepository) MarkFailed(ctx context.Context, id uuid.UUID, errorMsg string) error {
	query := `
		UPDATE outbox_events
		SET retry_count = retry_count + 1, last_error = $2
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, errorMsg)
	return err
}
