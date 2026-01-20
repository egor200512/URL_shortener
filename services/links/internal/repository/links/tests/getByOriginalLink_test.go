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

func TestGetByOriginalLink(t *testing.T) {
	pool := SetupLinksTestDB(t)
	repo := links.NewLinksRepo(context.Background(), pool)
	ctx := context.Background()

	tests := []struct {
		name         string
		setupLink    bool
		originalLink string
		wantLink     bool
	}{
		{
			name:         "link exists",
			setupLink:    true,
			originalLink: "example.com/path",
			wantLink:     true,
		},
		{
			name:         "link not exists",
			setupLink:    false,
			originalLink: "missing.com/path",
			wantLink:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pool.Exec(ctx, `TRUNCATE TABLE links.short_links`)
			require.NoError(t, err)

			if tt.setupLink {
				_, err := pool.Exec(ctx, `
                    INSERT INTO links.short_links (id, user_id, short_link, original_link_host, original_link, created_at)
                    VALUES ($1, $2, $3, $4, $5, $6)`,
					uuid.New(), uuid.New(), "abc123", "example.com", tt.originalLink, time.Now().UTC())
				require.NoError(t, err)
			}

			link, err := repo.GetByOriginalLink(ctx, tt.originalLink)

			assert.NoError(t, err)
			if tt.wantLink {
				assert.NotNil(t, link)
				assert.Equal(t, tt.originalLink, link.OriginalLink)
				assert.Equal(t, "example.com", link.OriginalLinkHost)
				assert.Equal(t, "abc123", link.ShortLink)
				assert.True(t, link.CreatedAt.Valid)
			} else {
				assert.Nil(t, link)
			}
		})
	}
}
