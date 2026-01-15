package tests

import (
	"context"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/services/auth/internal/models"
	"github.com/egor200512/URL_shortener/services/auth/internal/repository/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckRegistration(t *testing.T) {
	pool := SetupAuthTestDB(t)
	repo := auth.NewAuthRepo(context.Background(), pool)
	ctx := context.Background()

	tests := []struct {
		name       string
		setupUser  bool
		email      string
		wantUser   bool
		wantErr    bool
		verifyUser func(*testing.T, *models.User)
	}{
		{
			name:      "User exists",
			setupUser: true,
			email:     "test@example.com",
			wantUser:  true,
			wantErr:   false,
			verifyUser: func(t *testing.T, user *models.User) {
				assert.Equal(t, "test@example.com", user.Email)
				assert.Equal(t, []byte("sodium123"), user.Salt)
				assert.Equal(t, []byte("argon2id_hash"), user.SaltPassHash)
				assert.True(t, user.CreatedAt.Valid)
			},
		},
		{
			name:      "User not exists",
			setupUser: false,
			email:     "unknown@example.com",
			wantUser:  false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очистка таблицы перед каждым тестом
			_, err := pool.Exec(ctx, `TRUNCATE TABLE auth.users`)
			require.NoError(t, err)

			if tt.setupUser {
				// Вставка тестового пользователя
				userID := uuid.New()
				createdAt := time.Now().UTC()
				_, err := pool.Exec(ctx, `
                    INSERT INTO auth.users (id, email, salt, salt_password_hash, created_at) 
                    VALUES ($1, $2, $3, $4, $5)`,
					userID, tt.email, []byte("sodium123"), []byte("argon2id_hash"), createdAt)
				require.NoError(t, err)
			}

			user, err := repo.CheckRegistration(ctx, tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				if tt.wantUser {
					assert.NotNil(t, user)
					if tt.verifyUser != nil {
						tt.verifyUser(t, user)
					}
				} else {
					assert.Nil(t, user)
				}
			}
		})
	}
}
