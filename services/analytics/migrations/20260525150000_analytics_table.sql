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

CREATE INDEX link_events_short_link_idx ON analytics.link_events (short_link);
CREATE INDEX link_events_user_id_idx ON analytics.link_events (user_id);
CREATE INDEX link_events_executed_at_idx ON analytics.link_events (executed_at);

-- +goose Down
DROP TABLE IF EXISTS analytics.link_events;
DROP SCHEMA IF EXISTS analytics CASCADE;
