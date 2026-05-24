package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
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

	_, err = pool.Exec(context.Background(), `DROP SCHEMA IF EXISTS analytics CASCADE;`)
	require.NoError(t, err)

	_, err = pool.Exec(context.Background(), readGooseUp(t, filepath.Join("..", "..", "..", "..", "migrations", "20260123090108_analytics_table.sql")))
	require.NoError(t, err)

	t.Cleanup(func() { pool.Close() })
	return pool
}

func readGooseUp(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	sql := string(data)
	start := strings.Index(sql, "-- +goose Up")
	require.NotEqual(t, -1, start)

	end := strings.Index(sql, "-- +goose Down")
	require.NotEqual(t, -1, end)
	require.Greater(t, end, start)

	return strings.TrimSpace(sql[start+len("-- +goose Up") : end])
}
