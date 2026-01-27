package analytics

import (
	"context"
	"testing"
	"time"

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

func TestHandleMessage_OK(t *testing.T) {
	repo := &fakeRepo{}
	svc := &analyticsService{repo: repo}
	ts := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	msg := &nats.Msg{
		Subject: "links.created",
		Data:    []byte(`{"user_id":"u1","short_link":"s1","original_link":"o1","executed_at":"2026-01-01T00:00:00Z"}`),
	}

	if err := svc.HandleMessage(context.Background(), msg); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repo.called {
		t.Fatalf("expected repo to be called")
	}
	if repo.last == nil {
		t.Fatalf("expected payload to be set")
	}
	if repo.last.EventType != "created" {
		t.Fatalf("expected event_type=created, got %s", repo.last.EventType)
	}
	if repo.last.UserID != "u1" || repo.last.ShortLink != "s1" || repo.last.OriginalLink != "o1" {
		t.Fatalf("unexpected payload: %#v", repo.last)
	}
	if !repo.last.ExecutedAt.Equal(ts) {
		t.Fatalf("expected executed_at=%s, got %s", ts, repo.last.ExecutedAt)
	}
}

func TestHandleMessage_InvalidJSON(t *testing.T) {
	repo := &fakeRepo{}
	svc := &analyticsService{repo: repo}

	msg := &nats.Msg{
		Subject: "links.created",
		Data:    []byte(`{"user_id":`),
	}

	if err := svc.HandleMessage(context.Background(), msg); err == nil {
		t.Fatalf("expected error")
	}
	if repo.called {
		t.Fatalf("repo should not be called on invalid json")
	}
}
