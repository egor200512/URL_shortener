package analytics

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/services/analytics/internal/mocks"
	"github.com/egor200512/URL_shortener/services/analytics/internal/models"
	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAnalyticsService_RecordEvent(t *testing.T) {
	t.Parallel()

	executedAt := time.Date(2026, 5, 26, 10, 0, 0, 0, time.UTC)
	req := &desc.RecordEventRequest{
		EventType:    "created",
		UserId:       "user-id",
		ShortLink:    "abc123",
		OriginalLink: "example.com/path",
		ExecutedAt:   timestamppb.New(executedAt),
	}

	tests := []struct {
		name       string
		req        *desc.RecordEventRequest
		repoErr    error
		wantErrSub string
	}{
		{name: "success", req: req},
		{name: "repo error", req: req, repoErr: errors.New("insert failed"), wantErrSub: "insert failed"},
		{name: "default executed at", req: &desc.RecordEventRequest{EventType: "fetched", UserId: "user-id", ShortLink: "abc123", OriginalLink: "example.com/path"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewIAnalyticsRepo(t)
			repo.On("InsertEvent", mock.Anything, mock.MatchedBy(func(event *models.LinkEvent) bool {
				if event == nil {
					return false
				}
				if event.EventType != tt.req.EventType ||
					event.UserID != tt.req.UserId ||
					event.ShortLink != tt.req.ShortLink ||
					event.OriginalLink != tt.req.OriginalLink {
					return false
				}
				if tt.req.ExecutedAt != nil {
					return event.ExecutedAt.Equal(tt.req.ExecutedAt.AsTime())
				}
				return !event.ExecutedAt.IsZero()
			})).Return(tt.repoErr).Once()

			svc := NewAnalyticsService(repo)
			err := svc.RecordEvent(context.Background(), tt.req)
			assertErrContains(t, err, tt.wantErrSub)
		})
	}
}

func TestAnalyticsService_GetEvents(t *testing.T) {
	t.Parallel()

	executedAt := time.Date(2026, 5, 26, 10, 0, 0, 0, time.UTC)
	events := []*models.LinkEvent{
		{
			ID:           "event-id",
			EventType:    "created",
			UserID:       "user-id",
			ShortLink:    "abc123",
			OriginalLink: "example.com/path",
			ExecutedAt:   executedAt,
		},
	}

	tests := []struct {
		name       string
		repoErr    error
		wantErrSub string
	}{
		{name: "success"},
		{name: "repo error", repoErr: errors.New("db failed"), wantErrSub: "db failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewIAnalyticsRepo(t)
			repo.On("GetEvents", mock.Anything, int32(10), int32(2)).Return(events, int32(1), tt.repoErr).Once()

			svc := NewAnalyticsService(repo)
			resp, err := svc.GetEvents(context.Background(), 10, 2)
			assertErrContains(t, err, tt.wantErrSub)

			if tt.wantErrSub == "" {
				if resp.TotalCount != 1 || len(resp.Events) != 1 {
					t.Fatalf("response = %#v, want one event and total 1", resp)
				}
				if resp.Events[0].Id != events[0].ID || resp.Events[0].ExecutedAt.AsTime() != executedAt {
					t.Fatalf("event = %#v, want %#v", resp.Events[0], events[0])
				}
			}
		})
	}
}

func TestAnalyticsService_GetLinkStats(t *testing.T) {
	t.Parallel()

	stats := &models.LinkStats{
		ShortLink:    "abc123",
		CreatedCount: 1,
		FetchedCount: 2,
		DeletedCount: 3,
		TotalCount:   6,
	}

	tests := []struct {
		name       string
		repoErr    error
		wantErrSub string
	}{
		{name: "success"},
		{name: "repo error", repoErr: errors.New("db failed"), wantErrSub: "db failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewIAnalyticsRepo(t)
			repo.On("GetLinkStats", mock.Anything, "abc123").Return(stats, tt.repoErr).Once()

			svc := NewAnalyticsService(repo)
			resp, err := svc.GetLinkStats(context.Background(), "abc123")
			assertErrContains(t, err, tt.wantErrSub)

			if tt.wantErrSub == "" && (resp.ShortLink != stats.ShortLink || resp.TotalCount != stats.TotalCount) {
				t.Fatalf("stats = %#v, want %#v", resp, stats)
			}
		})
	}
}

func assertErrContains(t *testing.T, err error, wantSub string) {
	t.Helper()

	if wantSub == "" {
		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		return
	}
	if err == nil {
		t.Fatalf("err = nil, want containing %q", wantSub)
	}
	if !strings.Contains(err.Error(), wantSub) {
		t.Fatalf("err = %q, want containing %q", err.Error(), wantSub)
	}
}
