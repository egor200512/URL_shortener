package analytics

import (
	"context"
	"testing"
	"time"

	a "github.com/egor200512/URL_shortener/services/analytics/internal/service/analytics"
	"github.com/egor200512/URL_shortener/shared/models"
	"github.com/nats-io/nats.go"
)

type fakeRepo struct {
	called bool
	last   *models.LinkEvent
	err    error
}

func (f *fakeRepo) InsertEvent(ctx context.Context, payload *models.LinkEvent) error {
	f.called = true
	f.last = payload
	return f.err
}

func (f *fakeRepo) GetEvents(ctx context.Context, limit, offset int32) ([]*models.LinkEvent, error) {
	return nil, f.err
}

func TestHandleMessage(t *testing.T) {
	ts := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		msg         *nats.Msg
		expectErr   bool
		expectCall  bool
		expectEvent *models.LinkEvent
	}{
		{
			name: "ok",
			msg: &nats.Msg{
				Subject: "links.created",
				Data:    []byte(`{"user_id":"u1","short_link":"s1","original_link":"o1","executed_at":"2026-01-01T00:00:00Z"}`),
			},
			expectCall: true,
			expectEvent: &models.LinkEvent{
				EventType:    "created",
				UserID:       "u1",
				ShortLink:    "s1",
				OriginalLink: "o1",
				ExecutedAt:   ts,
			},
		},
		{
			name: "invalid json",
			msg: &nats.Msg{
				Subject: "links.created",
				Data:    []byte(`{"user_id":`),
			},
			expectErr:  true,
			expectCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}
			svc := a.NewAnalyticsService(repo)

			err := svc.HandleMessage(context.Background(), tt.msg)
			if tt.expectErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if repo.called != tt.expectCall {
				t.Fatalf("repo called=%t, want %t", repo.called, tt.expectCall)
			}

			if tt.expectEvent == nil {
				if repo.last != nil {
					t.Fatalf("expected no payload, got %#v", repo.last)
				}
				return
			}

			if repo.last == nil {
				t.Fatalf("expected payload to be set")
			}
			if repo.last.EventType != tt.expectEvent.EventType {
				t.Fatalf("event_type=%s, want %s", repo.last.EventType, tt.expectEvent.EventType)
			}
			if repo.last.UserID != tt.expectEvent.UserID || repo.last.ShortLink != tt.expectEvent.ShortLink || repo.last.OriginalLink != tt.expectEvent.OriginalLink {
				t.Fatalf("payload mismatch: got %#v, want %#v", repo.last, tt.expectEvent)
			}
			if !repo.last.ExecutedAt.Equal(tt.expectEvent.ExecutedAt) {
				t.Fatalf("executed_at=%s, want %s", repo.last.ExecutedAt, tt.expectEvent.ExecutedAt)
			}
		})
	}
}
