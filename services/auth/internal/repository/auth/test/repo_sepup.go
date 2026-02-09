package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

func SetupAuthTestDB(t *testing.T) *pgxpool.Pool {
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

	_, err = pool.Exec(context.Background(), `
        DROP SCHEMA IF EXISTS auth CASCADE;
        CREATE SCHEMA IF NOT EXISTS auth;

        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

        CREATE TABLE auth.users (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            email VARCHAR(100) UNIQUE NOT NULL,
            salt BYTEA NOT NULL,
            salt_password_hash BYTEA NOT NULL,
            created_at TIMESTAMP DEFAULT NOW()
        );
    `)
	require.NoError(t, err)

	t.Cleanup(func() { pool.Close() })
	return pool
}
