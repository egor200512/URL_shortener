-- +goose Up
CREATE SCHEMA IF NOT EXISTS links;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE links.short_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url_host TEXT NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS links.short_links;
DROP SCHEMA IF EXISTS links CASCADE;
