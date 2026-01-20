package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

func SetupLinksTestDB(t *testing.T) *pgxpool.Pool {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_AUTH_PORT_TESTS"),
		os.Getenv("PG_NAME"),
	)

	pool, err := pgxpool.Connect(context.Background(), connStr)
	require.NoError(t, err)

	err = pool.Ping(context.Background())
	require.NoError(t, err)

	_, err = pool.Exec(context.Background(), `
        DROP SCHEMA IF EXISTS links CASCADE;
        CREATE SCHEMA IF NOT EXISTS links;

        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

        CREATE TABLE links.short_links (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            user_id UUID NOT NULL,
            short_link VARCHAR(10) UNIQUE NOT NULL,
            original_link_host TEXT NOT NULL,
            original_link TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT NOW()
        );
    `)
	require.NoError(t, err)

	t.Cleanup(func() { pool.Close() })
	return pool
}
