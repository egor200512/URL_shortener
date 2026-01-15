package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/egor200512/URL_shortener/services/auth/internal/models"
	"github.com/egor200512/URL_shortener/services/auth/internal/repository/auth"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertUser(t *testing.T) {
	pool := SetupAuthTestDB(t)
	repo := auth.NewAuthRepo(context.Background(), pool)
	ctx := context.Background()

	longEmail := strings.Repeat("a", 101) // 101 символов

	tests := []struct {
		name            string
		regRequest      *models.RegisterRequest
		setupDB         func(*pgxpool.Pool)
		wantErr         bool
		wantErrContains string
		verifyWithCheck func(*testing.T, *models.RegisterRequest)
	}{
		{
			name: "successful insert",
			regRequest: &models.RegisterRequest{
				Email:        "success@example.com",
				Salt:         []byte("success_salt"),
				SaltPassHash: []byte("success_hash"),
			},
			wantErr: false,
			verifyWithCheck: func(t *testing.T, req *models.RegisterRequest) {
				user, err := repo.CheckRegistration(ctx, req.Email)
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, req.Email, user.Email)
				assert.Equal(t, req.Salt, user.Salt)
				assert.Equal(t, req.SaltPassHash, user.SaltPassHash)
			},
		},
		{
			name: "duplicate email",
			regRequest: &models.RegisterRequest{
				Email:        "duplicate@example.com",
				Salt:         []byte("dup_salt"),
				SaltPassHash: []byte("dup_hash"),
			},
			setupDB: func(pool *pgxpool.Pool) {
				_, err := pool.Exec(context.Background(),
					`INSERT INTO auth.users (id, email, salt, salt_password_hash) 
                     VALUES (uuid_generate_v4(), $1, $2, $3)`,
					"duplicate@example.com", []byte("existing"), []byte("existing_hash"))
				require.NoError(t, err)
			},
			wantErr:         true,
			wantErrContains: "duplicate key value violates unique constraint",
		},
		{
			name: "email too long",
			regRequest: &models.RegisterRequest{
				Email:        longEmail,
				Salt:         []byte("long_salt"),
				SaltPassHash: []byte("hash"),
			},
			wantErr:         true,
			wantErrContains: "value too long for type character varying(100)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очистка таблицы перед каждым тестом
			_, err := pool.Exec(ctx, `TRUNCATE TABLE auth.users`)
			require.NoError(t, err)

			if tt.setupDB != nil {
				tt.setupDB(pool)
			}

			err = repo.InsertUser(ctx, tt.regRequest)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrContains != "" {
					assert.Contains(t, err.Error(), tt.wantErrContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.verifyWithCheck != nil {
					tt.verifyWithCheck(t, tt.regRequest)
				}
			}
		})
	}
}
