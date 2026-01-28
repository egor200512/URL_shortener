package tests

import (
	"context"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/services/analytics/internal/repository/analytics"
)

func TestGetEvents(t *testing.T) {
	pool := SetupAnalyticsTestDB(t)
	repo := analytics.NewAnalyticsRepo(context.Background(), pool)
	ctx := context.Background()

	// Seed ordered data
	_, err := pool.Exec(ctx, `TRUNCATE TABLE analytics.link_events`)
	if err != nil {
		t.Fatalf("truncate failed: %v", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO analytics.link_events (id, event_type, user_id, short_link, original_link, executed_at)
		VALUES
		(uuid_generate_v4(), 'created', '11111111-1111-1111-1111-111111111111', 's1', 'o1', $1),
		(uuid_generate_v4(), 'fetched', '22222222-2222-2222-2222-222222222222', 's2', 'o2', $2),
		(uuid_generate_v4(), 'deleted', '33333333-3333-3333-3333-333333333333', 's3', 'o3', $3)
	`, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	events, err := repo.GetEvents(ctx, 2, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].ShortLink != "s3" || events[1].ShortLink != "s2" {
		t.Fatalf("unexpected order: %v, %v", events[0].ShortLink, events[1].ShortLink)
	}

	events, err = repo.GetEvents(ctx, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error (offset): %v", err)
	}
	if len(events) != 1 || events[0].ShortLink != "s1" {
		t.Fatalf("offset result mismatch, got %+v", events)
	}
}
