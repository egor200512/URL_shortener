package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

func SetupAnalyticsTestDB(t *testing.T) *pgxpool.Pool {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_TESTS_HOST"),
		os.Getenv("PG_PORT_TESTS"),
		os.Getenv("PG_NAME"),
	)

	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		t.Fatalf("failed to connect test db: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("failed to ping test db: %v", err)
	}

	_, err = pool.Exec(context.Background(), `
		DROP SCHEMA IF EXISTS analytics CASCADE;
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
	`)
	if err != nil {
		t.Fatalf("failed to setup schema: %v", err)
	}

	t.Cleanup(func() { pool.Close() })
	return pool
}
