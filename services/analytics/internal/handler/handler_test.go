package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/egor200512/URL_shortener/services/analytics/internal/mocks"
	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAnalyticsHandler_RecordEvent(t *testing.T) {
	t.Parallel()

	validReq := &desc.RecordEventRequest{
		EventType:    "created",
		UserId:       "user-id",
		ShortLink:    "abc123",
		OriginalLink: "example.com/path",
		ExecutedAt:   timestamppb.Now(),
	}

	tests := []struct {
		name     string
		req      *desc.RecordEventRequest
		svcErr   error
		wantCode codes.Code
		wantCall bool
	}{
		{name: "success", req: validReq, wantCode: codes.OK, wantCall: true},
		{name: "empty event type", req: cloneRecordEvent(validReq, func(req *desc.RecordEventRequest) { req.EventType = "" }), wantCode: codes.InvalidArgument},
		{name: "empty user id", req: cloneRecordEvent(validReq, func(req *desc.RecordEventRequest) { req.UserId = "" }), wantCode: codes.InvalidArgument},
		{name: "empty short link", req: cloneRecordEvent(validReq, func(req *desc.RecordEventRequest) { req.ShortLink = "" }), wantCode: codes.InvalidArgument},
		{name: "empty original link", req: cloneRecordEvent(validReq, func(req *desc.RecordEventRequest) { req.OriginalLink = "" }), wantCode: codes.InvalidArgument},
		{name: "service error", req: validReq, svcErr: errors.New("db failed"), wantCode: codes.Internal, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewIAnalyticsService(t)
			if tt.wantCall {
				svc.On("RecordEvent", mock.Anything, tt.req).Return(tt.svcErr).Once()
			}

			h := NewAnalyticsRouter(svc)
			resp, err := h.RecordEvent(context.Background(), tt.req)
			assertCode(t, err, tt.wantCode)

			if tt.wantCode == codes.OK && resp == nil {
				t.Fatal("expected response")
			}
		})
	}
}

func TestAnalyticsHandler_GetEvents(t *testing.T) {
	t.Parallel()

	resp := &desc.GetEventsResponse{TotalCount: 1}

	tests := []struct {
		name     string
		req      *desc.GetEventsRequest
		svcErr   error
		wantCode codes.Code
		wantCall bool
	}{
		{name: "success", req: &desc.GetEventsRequest{Limit: 10, Offset: 0}, wantCode: codes.OK, wantCall: true},
		{name: "bad limit", req: &desc.GetEventsRequest{Limit: 0}, wantCode: codes.InvalidArgument},
		{name: "bad offset", req: &desc.GetEventsRequest{Limit: 10, Offset: -1}, wantCode: codes.InvalidArgument},
		{name: "service error", req: &desc.GetEventsRequest{Limit: 10, Offset: 0}, svcErr: errors.New("db failed"), wantCode: codes.Internal, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewIAnalyticsService(t)
			if tt.wantCall {
				svc.On("GetEvents", mock.Anything, tt.req.Limit, tt.req.Offset).Return(resp, tt.svcErr).Once()
			}

			h := NewAnalyticsRouter(svc)
			got, err := h.GetEvents(context.Background(), tt.req)
			assertCode(t, err, tt.wantCode)

			if tt.wantCode == codes.OK && got != resp {
				t.Fatalf("response = %#v, want %#v", got, resp)
			}
		})
	}
}

func TestAnalyticsHandler_GetLinkStats(t *testing.T) {
	t.Parallel()

	resp := &desc.GetLinkStatsResponse{ShortLink: "abc123", TotalCount: 3}

	tests := []struct {
		name     string
		req      *desc.GetLinkStatsRequest
		svcErr   error
		wantCode codes.Code
		wantCall bool
	}{
		{name: "success", req: &desc.GetLinkStatsRequest{ShortLink: "abc123"}, wantCode: codes.OK, wantCall: true},
		{name: "empty short link", req: &desc.GetLinkStatsRequest{}, wantCode: codes.InvalidArgument},
		{name: "service error", req: &desc.GetLinkStatsRequest{ShortLink: "abc123"}, svcErr: errors.New("db failed"), wantCode: codes.Internal, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewIAnalyticsService(t)
			if tt.wantCall {
				svc.On("GetLinkStats", mock.Anything, tt.req.ShortLink).Return(resp, tt.svcErr).Once()
			}

			h := NewAnalyticsRouter(svc)
			got, err := h.GetLinkStats(context.Background(), tt.req)
			assertCode(t, err, tt.wantCode)

			if tt.wantCode == codes.OK && got != resp {
				t.Fatalf("response = %#v, want %#v", got, resp)
			}
		})
	}
}

func cloneRecordEvent(req *desc.RecordEventRequest, mutate func(*desc.RecordEventRequest)) *desc.RecordEventRequest {
	cp := *req
	mutate(&cp)
	return &cp
}

func assertCode(t *testing.T, err error, want codes.Code) {
	t.Helper()

	if got := status.Code(err); got != want {
		t.Fatalf("status code = %s, want %s; err = %v", got, want, err)
	}
}
