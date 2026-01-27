package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/services/analytics/internal/repository/analytics"
	"github.com/egor200512/URL_shortener/shared/models"
)

func TestInsertEvent(t *testing.T) {
	pool := SetupAnalyticsTestDB(t)
	repo := analytics.NewAnalyticsRepo(context.Background(), pool)
	ctx := context.Background()

	tests := []struct {
		name            string
		payload         *models.LinkEvent
		wantErr         bool
		wantErrContains string
	}{
		{
			name: "success",
			payload: &models.LinkEvent{
				EventType:    "created",
				UserID:       "11111111-1111-1111-1111-111111111111",
				ShortLink:    "abc123",
				OriginalLink: "example.com/path",
				ExecutedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid event type",
			payload: &models.LinkEvent{
				EventType:    "bad",
				UserID:       "11111111-1111-1111-1111-111111111111",
				ShortLink:    "abc124",
				OriginalLink: "example.com/path",
				ExecutedAt:   time.Now(),
			},
			wantErr:         true,
			wantErrContains: "violates check constraint",
		},
		{
			name: "invalid user id",
			payload: &models.LinkEvent{
				EventType:    "created",
				UserID:       "not-a-uuid",
				ShortLink:    "abc125",
				OriginalLink: "example.com/path",
				ExecutedAt:   time.Now(),
			},
			wantErr:         true,
			wantErrContains: "invalid input syntax for type uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pool.Exec(ctx, `TRUNCATE TABLE analytics.link_events`)
			if err != nil {
				t.Fatalf("failed to truncate: %v", err)
			}

			err = repo.InsertEvent(ctx, tt.payload)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				if tt.wantErrContains != "" && err != nil && !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Fatalf("expected error to contain %q, got %v", tt.wantErrContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var count int
			err = pool.QueryRow(ctx, `
				SELECT COUNT(*) FROM analytics.link_events
				WHERE short_link = $1 AND original_link = $2`,
				tt.payload.ShortLink, tt.payload.OriginalLink,
			).Scan(&count)
			if err != nil {
				t.Fatalf("failed to query: %v", err)
			}
			if count != 1 {
				t.Fatalf("expected 1 row, got %d", count)
			}
		})
	}
}
