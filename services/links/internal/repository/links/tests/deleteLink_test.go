package tests

import (
	"context"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/services/links/internal/repository/links"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteLink(t *testing.T) {
	pool := SetupLinksTestDB(t)
	repo := links.NewLinksRepo(context.Background(), pool)
	ctx := context.Background()

	tests := []struct {
		name            string
		setupLink       func() (shortLink string, userID uuid.UUID)
		deleteShortLink string
		deleteUserID    uuid.UUID
		verify          func(t *testing.T, shortLink string)
	}{
		{
			name: "successful delete",
			setupLink: func() (string, uuid.UUID) {
				userID := uuid.New()
				shortLink := "abc123"
				_, err := pool.Exec(ctx, `
                    INSERT INTO links.short_links (id, user_id, short_link, original_link_host, original_link, created_at)
                    VALUES ($1, $2, $3, $4, $5, $6)`,
					uuid.New(), userID, shortLink, "example.com", "example.com/path", time.Now().UTC())
				require.NoError(t, err)
				return shortLink, userID
			},
			deleteShortLink: "abc123",
			verify: func(t *testing.T, shortLink string) {
				var count int
				err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM links.short_links WHERE short_link = $1`, shortLink).Scan(&count)
				require.NoError(t, err)
				assert.Equal(t, 0, count)
			},
		},
		{
			name: "wrong userID - link remains",
			setupLink: func() (string, uuid.UUID) {
				userID := uuid.New()
				shortLink := "abc123"
				_, err := pool.Exec(ctx, `
                    INSERT INTO links.short_links (id, user_id, short_link, original_link_host, original_link, created_at)
                    VALUES ($1, $2, $3, $4, $5, $6)`,
					uuid.New(), userID, shortLink, "example.com", "example.com/path", time.Now().UTC())
				require.NoError(t, err)
				return shortLink, userID
			},
			deleteShortLink: "abc123",
			deleteUserID:    uuid.New(),
			verify: func(t *testing.T, shortLink string) {
				var count int
				err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM links.short_links WHERE short_link = $1`, shortLink).Scan(&count)
				require.NoError(t, err)
				assert.Equal(t, 1, count)
			},
		},
		{
			name: "non-existent link - no error",
			setupLink: func() (string, uuid.UUID) {
				return "missing", uuid.New()
			},
			deleteShortLink: "missing",
			verify: func(t *testing.T, shortLink string) {
				var count int
				err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM links.short_links WHERE short_link = $1`, shortLink).Scan(&count)
				require.NoError(t, err)
				assert.Equal(t, 0, count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pool.Exec(ctx, `TRUNCATE TABLE links.short_links`)
			require.NoError(t, err)

			shortLink, userID := tt.setupLink()
			if tt.deleteUserID == uuid.Nil {
				tt.deleteUserID = userID
			}

			err = repo.DeleteLink(ctx, tt.deleteShortLink, tt.deleteUserID.String())
			assert.NoError(t, err)

			if tt.verify != nil {
				tt.verify(t, shortLink)
			}
		})
	}
}
