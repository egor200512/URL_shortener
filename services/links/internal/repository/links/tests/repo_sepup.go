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

func SetupLinksTestDB(t *testing.T) *pgxpool.Pool {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_TESTS_HOST"),
		os.Getenv("PG_PORT_TESTS"),
		os.Getenv("PG_NAME"),
	)

	pool, err := pgxpool.Connect(context.Background(), connStr)
	require.NoError(t, err)

	err = pool.Ping(context.Background())
	require.NoError(t, err)

	_, err = pool.Exec(context.Background(), `DROP SCHEMA IF EXISTS links CASCADE;`)
	require.NoError(t, err)

	_, err = pool.Exec(context.Background(), readGooseUp(t, filepath.Join("..", "..", "..", "..", "migrations", "20260117134530_links_table.sql")))
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
