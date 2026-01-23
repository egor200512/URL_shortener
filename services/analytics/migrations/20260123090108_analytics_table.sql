-- +goose Up
CREATE SCHEMA IF NOT EXISTS analytics;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE analytics.link_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type TEXT NOT NULL CHECK (event_type IN ('created', 'fetched', 'deleted')),
    user_id UUID NOT NULL,
    short_link VARCHAR(10) NOT NULL,
    original_link TEXT NOT NULL,
    executed_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS analytics.link_events;
DROP SCHEMA IF EXISTS analytics CASCADE;
