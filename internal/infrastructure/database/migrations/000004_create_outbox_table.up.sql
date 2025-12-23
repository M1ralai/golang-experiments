-- Outbox pattern for reliable event publishing
CREATE TABLE IF NOT EXISTS outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP,
    retry_count INTEGER DEFAULT 0,
    last_error TEXT,
    INDEX idx_outbox_unprocessed (processed_at, created_at) WHERE processed_at IS NULL
);

COMMENT ON TABLE outbox_events IS 'Transactional outbox for reliable event publishing';
COMMENT ON COLUMN outbox_events.aggregate_type IS 'Entity type (e.g., task, user)';
COMMENT ON COLUMN outbox_events.aggregate_id IS 'Entity ID that triggered the event';
COMMENT ON COLUMN outbox_events.event_type IS 'Event name (e.g., TaskAssigned, UserCreated)';
COMMENT ON COLUMN outbox_events.payload IS 'Event payload as JSON';
COMMENT ON COLUMN outbox_events.processed_at IS 'NULL = pending, timestamp = processed';
