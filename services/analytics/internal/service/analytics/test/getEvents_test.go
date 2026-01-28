package analytics

import (
	"context"
	"errors"
	"testing"
	"time"

	a "github.com/egor200512/URL_shortener/services/analytics/internal/service/analytics"
	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"github.com/egor200512/URL_shortener/shared/models"
)

type fakeRepoEvents struct {
	items []*models.LinkEvent
	err   error
}

func (f *fakeRepoEvents) InsertEvent(ctx context.Context, payload *models.LinkEvent) error {
	return nil
}

func (f *fakeRepoEvents) GetEvents(ctx context.Context, limit, offset int32) ([]*models.LinkEvent, error) {
	return f.items, f.err
}

func TestServiceGetEvents(t *testing.T) {
	ts := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		repoItems []*models.LinkEvent
		repoErr   error
		wantErr   bool
	}{
		{
			name: "success",
			repoItems: []*models.LinkEvent{
				{
					EventType:    "created",
					UserID:       "u1",
					ShortLink:    "s1",
					OriginalLink: "o1",
					ExecutedAt:   ts,
				},
				{
					EventType:    "FETCHED",
					UserID:       "u2",
					ShortLink:    "s2",
					OriginalLink: "o2",
					ExecutedAt:   ts.Add(time.Hour),
				},
			},
		},
		{
			name:    "repo error",
			repoErr: errors.New("boom"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepoEvents{items: tt.repoItems, err: tt.repoErr}
			svc := a.NewAnalyticsService(repo)

			resp, err := svc.GetEvents(context.Background(), &desc.GetEventsRequest{Limit: 10, Offset: 0})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(resp.Events) != len(tt.repoItems) {
				t.Fatalf("expected %d events, got %d", len(tt.repoItems), len(resp.Events))
			}

			// spot-check first item mapping
			got := resp.Events[0]
			want := tt.repoItems[0]
			if got.GetUserId() != want.UserID || got.GetShortLink() != want.ShortLink || got.GetOriginalLink() != want.OriginalLink {
				t.Fatalf("payload mismatch: got %+v, want %+v", got, want)
			}
			if got.GetType() != desc.LinkEventType_LINK_EVENT_TYPE_CREATED {
				t.Fatalf("expected created enum, got %v", got.GetType())
			}
			if got.GetOccurredAt().AsTime() != want.ExecutedAt {
				t.Fatalf("timestamp mismatch: got %s, want %s", got.GetOccurredAt().AsTime(), want.ExecutedAt)
			}

			if len(resp.Events) > 1 && resp.Events[1].GetType() != desc.LinkEventType_LINK_EVENT_TYPE_FETCHED {
				t.Fatalf("second event type mapping failed: %v", resp.Events[1].GetType())
			}
		})
	}
}
