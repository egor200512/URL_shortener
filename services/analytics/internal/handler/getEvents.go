package handler

import (
	"context"
	"fmt"
	"log/slog"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *AnalyticsHandler) GetEvents(ctx context.Context, in *desc.GetEventsRequest) (*desc.GetEventsResponse, error) {
	slog.Info("get events endpoint called", "limit", in.Limit, "offset", in.Offset)

	if in.Limit <= 0 {
		slog.Warn("get events validation failed", "limit", in.Limit, "error", "limit must be positive")
		return nil, status.Error(codes.InvalidArgument, "limit must be positive")
	}
	if in.Offset < 0 {
		slog.Warn("get events validation failed", "offset", in.Offset, "error", "offset must be non-negative")
		return nil, status.Error(codes.InvalidArgument, "offset must be non-negative")
	}

	events, err := handler.analyticsService.GetEvents(ctx, in.Limit, in.Offset)
	if err != nil {
		slog.Error("get events failed", "limit", in.Limit, "offset", in.Offset, "error", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get events: %s", err.Error()))
	}

	slog.Info("get events completed", "limit", in.Limit, "offset", in.Offset, "count", len(events.Events), "total_count", events.TotalCount)
	return events, nil
}
