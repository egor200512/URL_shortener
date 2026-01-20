package tests

import (
	"context"
	"testing"

	"github.com/egor200512/URL_shortener/services/links/internal/repository/links"
	"github.com/egor200512/URL_shortener/services/links/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertLink(t *testing.T) {
	pool := SetupLinksTestDB(t)
	repo := links.NewLinksRepo(context.Background(), pool)
	ctx := context.Background()

	tests := []struct {
		name            string
		req             *models.CreateLinkReq
		setupDB         func()
		wantErr         bool
		wantErrContains string
		verify          func(t *testing.T, req *models.CreateLinkReq)
	}{
		{
			name: "successful insert",
			req: &models.CreateLinkReq{
				UserID:           uuid.New().String(),
				ShortLink:        "abc123",
				OriginalLinkHost: "example.com",
				OriginalLink:     "example.com/path",
			},
			wantErr: false,
			verify: func(t *testing.T, req *models.CreateLinkReq) {
				var count int
				err := pool.QueryRow(ctx, `
                    SELECT COUNT(*) FROM links.short_links
                    WHERE short_link = $1 AND original_link = $2`,
					req.ShortLink, req.OriginalLink).Scan(&count)
				require.NoError(t, err)
				assert.Equal(t, 1, count)
			},
		},
		{
			name: "duplicate short link",
			req: &models.CreateLinkReq{
				UserID:           uuid.New().String(),
				ShortLink:        "dup123",
				OriginalLinkHost: "example.com",
				OriginalLink:     "example.com/dup",
			},
			setupDB: func() {
				_, err := pool.Exec(ctx, `
                    INSERT INTO links.short_links (id, user_id, short_link, original_link_host, original_link)
                    VALUES ($1, $2, $3, $4, $5)`,
					uuid.New(), uuid.New(), "dup123", "example.com", "example.com/exists")
				require.NoError(t, err)
			},
			wantErr:         true,
			wantErrContains: "duplicate key value violates unique constraint",
		},
		{
			name: "invalid userID",
			req: &models.CreateLinkReq{
				UserID:           "not-a-uuid",
				ShortLink:        "baduid",
				OriginalLinkHost: "example.com",
				OriginalLink:     "example.com/bad",
			},
			wantErr:         true,
			wantErrContains: "invalid input syntax for type uuid",
		},
		{
			name: "short link too long",
			req: &models.CreateLinkReq{
				UserID:           uuid.New().String(),
				ShortLink:        "toolong12345",
				OriginalLinkHost: "example.com",
				OriginalLink:     "example.com/long",
			},
			wantErr:         true,
			wantErrContains: "value too long for type character varying(10)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pool.Exec(ctx, `TRUNCATE TABLE links.short_links`)
			require.NoError(t, err)

			if tt.setupDB != nil {
				tt.setupDB()
			}

			err = repo.InsertLink(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrContains != "" {
					assert.Contains(t, err.Error(), tt.wantErrContains)
				}
				return
			}

			assert.NoError(t, err)
			if tt.verify != nil {
				tt.verify(t, tt.req)
			}
		})
	}
}
