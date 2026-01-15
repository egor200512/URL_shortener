-- +goose Up
CREATE SCHEMA IF NOT EXISTS auth;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE auth.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(100) UNIQUE NOT NULL,
    salt BYTEA NOT NULL,
    salt_password_hash BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS auth.users;
DROP SCHEMA IF EXISTS auth CASCADE;
