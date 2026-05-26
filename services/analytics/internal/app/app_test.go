package app

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/services/analytics/internal/mocks"
	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"github.com/egor200512/URL_shortener/shared/pkg/events"
	"github.com/stretchr/testify/mock"
)

func TestApp_handleLinkEvent(t *testing.T) {
	t.Parallel()

	executedAt := time.Date(2026, 5, 26, 10, 0, 0, 0, time.UTC)
	validEvent := events.LinkEvent{
		EventID:      "event-id",
		EventType:    events.LinkCreated,
		UserID:       "user-id",
		ShortLink:    "abc123",
		OriginalLink: "example.com/path",
		ExecutedAt:   executedAt,
	}

	tests := []struct {
		name       string
		payload    []byte
		wantCall   bool
		wantErrSub string
	}{
		{name: "bad json", payload: []byte("{bad"), wantErrSub: "invalid character"},
		{name: "empty event type", payload: mustMarshalEvent(t, cloneLinkEvent(validEvent, func(event *events.LinkEvent) { event.EventType = "" })), wantErrSub: "event type is required"},
		{name: "empty user id", payload: mustMarshalEvent(t, cloneLinkEvent(validEvent, func(event *events.LinkEvent) { event.UserID = "" })), wantErrSub: "user id is required"},
		{name: "empty short link", payload: mustMarshalEvent(t, cloneLinkEvent(validEvent, func(event *events.LinkEvent) { event.ShortLink = "" })), wantErrSub: "short link is required"},
		{name: "empty original link", payload: mustMarshalEvent(t, cloneLinkEvent(validEvent, func(event *events.LinkEvent) { event.OriginalLink = "" })), wantErrSub: "original link is required"},
		{name: "success", payload: mustMarshalEvent(t, validEvent), wantCall: true},
		{name: "record event error", payload: mustMarshalEvent(t, validEvent), wantCall: true, wantErrSub: "record failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewIAnalyticsService(t)
			if tt.wantCall {
				svc.On("RecordEvent", mock.Anything, mock.MatchedBy(func(req *desc.RecordEventRequest) bool {
					return req.EventType == validEvent.EventType &&
						req.UserId == validEvent.UserID &&
						req.ShortLink == validEvent.ShortLink &&
						req.OriginalLink == validEvent.OriginalLink &&
						req.ExecutedAt.AsTime().Equal(validEvent.ExecutedAt)
				})).Return(errorForTest(tt.wantErrSub)).Once()
			}

			a := &App{AnalyticsService: svc}
			err := a.handleLinkEvent(context.Background(), tt.payload)
			assertErrContains(t, err, tt.wantErrSub)
		})
	}
}

func cloneLinkEvent(event events.LinkEvent, mutate func(*events.LinkEvent)) events.LinkEvent {
	mutate(&event)
	return event
}

func mustMarshalEvent(t *testing.T, event events.LinkEvent) []byte {
	t.Helper()

	payload, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	return payload
}

func errorForTest(wantErrSub string) error {
	if wantErrSub == "record failed" {
		return errors.New(wantErrSub)
	}
	return nil
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
