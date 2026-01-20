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

func TestGetUserLinks(t *testing.T) {
	pool := SetupLinksTestDB(t)
	repo := links.NewLinksRepo(context.Background(), pool)
	ctx := context.Background()

	userID := uuid.New()
	otherUserID := uuid.New()
	now := time.Now().UTC()

	insertLink := func(uID uuid.UUID, shortLink, originalLink string, createdAt time.Time) {
		_, err := pool.Exec(ctx, `
            INSERT INTO links.short_links (id, user_id, short_link, original_link_host, original_link, created_at)
            VALUES ($1, $2, $3, $4, $5, $6)`,
			uuid.New(), uID, shortLink, "example.com", originalLink, createdAt)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		setupDB       func()
		limit         int32
		offset        int32
		expectedCount int
		expectedTotal int32
		verify        func(t *testing.T, linksCount int)
	}{
		{
			name: "success - returns paged links",
			setupDB: func() {
				insertLink(userID, "link1", "example.com/1", now.Add(-2*time.Hour))
				insertLink(userID, "link2", "example.com/2", now.Add(-1*time.Hour))
				insertLink(userID, "link3", "example.com/3", now)
				insertLink(otherUserID, "other1", "example.com/other", now)
			},
			limit:         2,
			offset:        0,
			expectedCount: 2,
			expectedTotal: 3,
		},
		{
			name: "success - offset works",
			setupDB: func() {
				insertLink(userID, "link1", "example.com/1", now.Add(-2*time.Hour))
				insertLink(userID, "link2", "example.com/2", now.Add(-1*time.Hour))
				insertLink(userID, "link3", "example.com/3", now)
			},
			limit:         2,
			offset:        2,
			expectedCount: 1,
			expectedTotal: 3,
		},
		{
			name: "success - no links",
			setupDB: func() {
				insertLink(otherUserID, "other1", "example.com/other", now)
			},
			limit:         5,
			offset:        0,
			expectedCount: 0,
			expectedTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pool.Exec(ctx, `TRUNCATE TABLE links.short_links`)
			require.NoError(t, err)

			if tt.setupDB != nil {
				tt.setupDB()
			}

			linksResult, total, err := repo.GetUserLinks(ctx, userID.String(), tt.limit, tt.offset)
			assert.NoError(t, err)
			assert.Len(t, linksResult, tt.expectedCount)
			assert.Equal(t, tt.expectedTotal, total)
		})
	}
}
