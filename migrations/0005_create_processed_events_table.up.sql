CREATE TABLE IF NOT EXISTS processed_events (
    event_type TEXT NOT NULL,
    entity_id  INTEGER NOT NULL,
    processed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_type, entity_id)
);
CREATE INDEX idx_processed_events_event_type ON processed_events(event_type);